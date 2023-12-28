package pumper

import (
	"encoding/json"
	"fmt"
	"runtime/debug"
	"strconv"
	"sync"

	"github.com/cralack/ChaosMetrics/server/global"
	"github.com/cralack/ChaosMetrics/server/model/riotmodel"
	"github.com/cralack/ChaosMetrics/server/service/fetcher"
	"github.com/cralack/ChaosMetrics/server/utils"
	"github.com/cralack/ChaosMetrics/server/utils/scheduler"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const testSize = 1
const apiToken = ""

type Pumper struct {
	logger      *zap.Logger
	db          *gorm.DB
	rdb         *redis.Client
	lock        *sync.Mutex
	fetcher     fetcher.Fetcher
	scheduler   scheduler.Scheduler
	entryMap    map[string]map[string]*riotmodel.LeagueEntryDTO // entryMap[Loc][SummonerID]
	sumnMap     map[string]map[string]*riotmodel.SummonerDTO    // sumnMap[Loc][SummonerID]
	matchMap    map[string]map[string]bool                      // matchMap[Loc][matchID]
	out         chan *DBResult
	entrieIdx   []uint
	summonerIdx []uint
	stgy        *Strategy
}

type DBResult struct {
	Type  string
	Brief string
	Data  interface{}
}

func NewPumper(opts ...Option) *Pumper {
	stgy := defaultStrategy
	for _, opt := range opts {
		opt(stgy)
	}

	if global.GVA_ENV == global.TEST_ENV {
		stgy.MaxMatchCount = 1
	}

	return &Pumper{
		logger:      global.GVA_LOG,
		db:          global.GVA_DB,
		rdb:         global.GVA_RDB,
		lock:        &sync.Mutex{},
		fetcher:     fetcher.NewBrowserFetcher(),
		scheduler:   scheduler.NewSchdule(),
		entrieIdx:   make([]uint, 16),
		summonerIdx: make([]uint, 16),
		out:         make(chan *DBResult),
		entryMap:    make(map[string]map[string]*riotmodel.LeagueEntryDTO),
		sumnMap:     make(map[string]map[string]*riotmodel.SummonerDTO),
		matchMap:    make(map[string]map[string]bool),

		stgy: stgy,
	}
}

func (p *Pumper) Schedule() {
	p.scheduler.Schedule()
}

func (p *Pumper) handleResult(exit chan struct{}) {
	for result := range p.out {
		switch result.Type {
		case "finish":
			p.logger.Info(fmt.Sprintf("all %s result store done", result.Brief))
			exit <- struct{}{}
			continue

		case "entry":
			entries := result.Data.([]*riotmodel.LeagueEntryDTO)
			if err := p.db.Save(entries).Error; err != nil {
				p.logger.Error("riot entry store failed", zap.Error(err))
			}

		case "summoners":
			summoners := result.Data.([]*riotmodel.SummonerDTO)
			if err := p.db.Save(summoners).Error; err != nil {
				p.logger.Error("riot summoner model store failed", zap.Error(err))
			}

		case "match":
			matches := result.Data.([]*riotmodel.MatchDB)
			if len(matches) == 0 {
				continue
			}
			if err := p.db.Create(matches).Error; err != nil {
				p.logger.Error(result.Brief+"'s match store failed", zap.Error(err))
			}
		}
	}
}

func (p *Pumper) StartEngine(exit chan struct{}) {
	go p.Schedule()
	// distrib part
	// go p.loadResource()
	// go p.watchResource()
	go p.fetch()
	go p.handleResult(exit)
}

func (p *Pumper) UpdateAll() {
	exit := make(chan struct{})
	p.StartEngine(exit)

	p.UpdateEntries(exit)
	p.UpdateSumoner(exit)
	p.UpdateMatch(exit)
}

// core func
func (p *Pumper) fetch() {
	var (
		cnt          int
		buff         []byte
		err          error
		list         *riotmodel.LeagueListDTO
		entries      []*riotmodel.LeagueEntryDTO
		matches      []*riotmodel.MatchDB
		curMatchList []string
		endTier      string
		endRank      string
	)
	endTier, endRank = ConvertRankToStr(p.stgy.TestEndMark[0], p.stgy.TestEndMark[1])
	// catch panic
	defer func() {
		if err := recover(); err != nil {
			p.logger.Panic("fetcher panic",
				zap.Any("err", err),
				zap.String("stack", string(debug.Stack())))
		}
	}()

	for {
		req := p.scheduler.Pull()
		if req.Data == nil {
			// send finish signal
			p.out <- &DBResult{
				Type:  "finish",
				Brief: req.Type,
				Data:  nil,
			}
			continue
		}

		// fetch and parse
		switch req.Type {
		case "bestEntry":
			data := req.Data.(*entryTask)
			// api: /lol/league/v4/{BEST}leagues/by-queue/{queue}
			if buff, err = p.fetcher.Get(fetcher.NewTask(
				fetcher.WithURL(req.URL),
				fetcher.WithToken(apiToken),
			)); err != nil || buff == nil {
				p.logger.Error(fmt.Sprintf("fetch %s %s failed", data.Tier, data.Rank),
					zap.Error(err))
				// fetch again
				if req.Retry < p.stgy.Retry {
					req.Retry++
					p.scheduler.Push(req)
				}
				continue
			}
			if err = json.Unmarshal(buff, &list); err != nil {
				p.logger.Error(fmt.Sprintf("unmarshal json to %s failed",
					"LeagueListDTO"), zap.Error(err))
			}
			entries = list.Entries
			if len(entries) == 0 {
				continue
			}
			for _, e := range entries {
				e.Tier = data.Tier
				e.Loc = req.Loc
			}
			// shrink size if test
			if global.GVA_ENV == global.TEST_ENV {
				entries = entries[:testSize]
			}

			list = nil
			p.logger.Info(fmt.Sprintf("all %d %s data fetch done",
				len(entries), data.Tier))
			p.handleEntries(entries, req.Loc)
			p.cacheEntries(entries, req.Loc)
			if data.Tier == endTier && data.Rank == endRank {
				p.out <- &DBResult{
					Type:  "finish",
					Brief: "entry",
					Data:  nil,
				}
				// *need release scheduler resource*
				continue
			}

		case "mortalEntry":
			data := req.Data.(*entryTask)
			for page := 1; ; page++ {
				// api: /lol/league/v4/entries/{queue}/{tier}/{division}
				url := fmt.Sprintf("%s?page=%s", req.URL, strconv.Itoa(page))
				if buff, err = p.fetcher.Get(fetcher.NewTask(
					fetcher.WithURL(url),
					fetcher.WithToken(apiToken),
				)); err != nil {
					p.logger.Error(fmt.Sprintf("fetch %s %s failed", data.Tier, data.Rank),
						zap.Error(err))
					if req.Retry < p.stgy.Retry {
						req.Retry++
						p.scheduler.Push(req)
					}
					continue
				}
				if err = json.Unmarshal(buff, &entries); err != nil {
					p.logger.Error(fmt.Sprintf("unmarshal json to %s failed",
						"LeagueEntryDTO"), zap.Error(err))
				} else {
					p.logger.Info(fmt.Sprintf("fetch %s %s page %d done", data.Tier, data.Rank, page))
				}
				for _, e := range entries {
					e.Loc = req.Loc
				}
				// shrink size if test
				if global.GVA_ENV == global.TEST_ENV {
					entries = entries[:testSize]
				}

				p.handleEntries(entries, req.Loc)
				p.cacheEntries(entries, req.Loc)

				// test
				if (global.GVA_ENV == global.TEST_ENV && page == testSize) || len(entries) == 0 {
					p.logger.Info(fmt.Sprintf("all %s %s data fetch done at page %d",
						data.Tier, data.Rank, page))
					break
				}
			}

			if data.Tier == endTier && data.Rank == endRank {
				p.out <- &DBResult{
					Type:  "finish",
					Brief: "entry",
					Data:  nil,
				}
				continue
			}

		case "summoner":
			data := req.Data.(*summonerTask)
			if buff, err = p.fetcher.Get(fetcher.NewTask(
				fetcher.WithURL(req.URL),
				fetcher.WithToken(apiToken),
			)); err != nil || buff == nil {
				p.logger.Error(fmt.Sprintf("fetch summonerID %s failed", data.summonerID), zap.Error(err))
				// fetch again
				if req.Retry < p.stgy.Retry {
					req.Retry++
					p.scheduler.Push(req)
				}
				continue
			}
			var sumn *riotmodel.SummonerDTO
			if err = json.Unmarshal(buff, &sumn); err != nil {
				p.logger.Error(fmt.Sprintf("unmarshal json to %s failed",
					"sumnDTO"), zap.Error(err))
			}
			p.handleSummoner(req.Loc, sumn)

		case "match":
			data := req.Data.(*matchTask)
			// get old & cur match list
			if buff, err = p.fetcher.Get(fetcher.NewTask(
				fetcher.WithURL(req.URL),
				fetcher.WithToken(apiToken),
			)); err != nil || buff == nil {
				p.logger.Error(fmt.Sprintf("fetch summoner %s's match list failed",
					data.sumn.Name), zap.Error(err))
				if req.Retry < p.stgy.Retry {
					req.Retry++
					p.scheduler.Push(req)
				}
				continue
			}
			if err = json.Unmarshal(buff, &curMatchList); err != nil {
				p.logger.Error(fmt.Sprintf("unmarshal json to %s failed",
					"match"), zap.Error(err))
			}
			// get old match list
			oldMatchList := make(map[string]struct{})
			for _, matchID := range utils.ConvertStrToSlice(data.sumn.Matches) {
				oldMatchList[matchID] = struct{}{}
			}
			// update summoner's match list
			summoner := data.sumn
			summoner.Matches = utils.ConvertSliceToStr(curMatchList)
			cnt++
			p.handleSummoner(req.Loc, summoner)
			// init val
			matches = make([]*riotmodel.MatchDB, 0, p.stgy.MaxMatchCount)
			loc := utils.ConverHostLoCode(req.Loc)
			host := utils.ConvertPlatformToHost(loc)
			// fetch match
			for _, matchID := range curMatchList {
				if _, has := p.matchMap[req.Loc][matchID]; has {
					continue
				}
				if _, has := oldMatchList[matchID]; has {
					continue
				} else {
					p.matchMap[req.Loc][matchID] = true
					if tmp := p.FetchMatchByID(req, host, matchID); tmp != nil {
						matches = append(matches, tmp)
					}
				}
			}
			p.logger.Debug(fmt.Sprintf("updating %s's match list @ %d,store %d matches",
				summoner.Name, cnt, len(matches)))
			if len(matches) == 0 {
				continue
			}
			p.handleMatches(matches, summoner.Name)
		}
	}
}

func getQueueString(que uint) string {
	switch que {
	case riotmodel.RANKED_SOLO_5x5:
		return "RANKED_SOLO_5x5"
	case riotmodel.RANKED_FLEX_SR:
		return "RANKED_FLEX_SR"
	case riotmodel.RANKED_FLEX_TT:
		return "RANKED_FLEX_TT"
	default:
		return ""
	}
}

func ConvertRankToStr(tier, div uint) (string, string) {
	var d string
	switch div {
	case 1:
		d = "I"
	case 2:
		d = "II"
	case 3:
		d = "III"
	case 4:
		d = "IV"
	}
	switch tier {
	case riotmodel.CHALLENGER:
		return "CHALLENGER", "I"
	case riotmodel.GRANDMASTER:
		return "GRANDMASTER", "I"
	case riotmodel.MASTER:
		return "MASTER", "I"
	case riotmodel.DIAMOND:
		return "DIAMOND", d
	case riotmodel.EMERALD:
		return "EMERALD", d
	case riotmodel.PLATINUM:
		return "PLATINUM", d
	case riotmodel.GOLD:
		return "GOLD", d
	case riotmodel.SILVER:
		return "SILVER", d
	case riotmodel.BRONZE:
		return "BRONZE", d
	case riotmodel.IRON:
		return "IRON", d
	}

	return "", ""
}

func ConvertStrToRank(tierStr, divStr string) (uint, uint) {
	var tier uint
	var div uint

	switch tierStr {
	case "CHALLENGER":
		tier = riotmodel.CHALLENGER
	case "GRANDMASTER":
		tier = riotmodel.GRANDMASTER
	case "MASTER":
		tier = riotmodel.MASTER
	case "DIAMOND":
		tier = riotmodel.DIAMOND
	case "EMERALD":
		tier = riotmodel.EMERALD
	case "PLATINUM":
		tier = riotmodel.PLATINUM
	case "GOLD":
		tier = riotmodel.GOLD
	case "SILVER":
		tier = riotmodel.SILVER
	case "BRONZE":
		tier = riotmodel.BRONZE
	case "IRON":
		tier = riotmodel.IRON
	default:
		return 0, 0
	}

	switch divStr {
	case "I":
		div = 1
	case "II":
		div = 2
	case "III":
		div = 3
	case "IV":
		div = 4
	default:
		return 0, 0
	}

	return tier, div
}
