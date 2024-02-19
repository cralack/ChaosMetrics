package pumper

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"sync"
	"time"

	"github.com/cralack/ChaosMetrics/server/internal/global"
	"github.com/cralack/ChaosMetrics/server/internal/service/fetcher"
	"github.com/cralack/ChaosMetrics/server/model/riotmodel"
	"github.com/cralack/ChaosMetrics/server/utils"
	"github.com/cralack/ChaosMetrics/server/utils/scheduler"
	"github.com/redis/go-redis/v9"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const testSize = 1

const (
	EntryTypeKey       = "entry"
	bestEntryTypeKey   = "best"
	mortalEntryTypeKey = "mortal"
	SummonerTypeKey    = "sum"
	MatchTypeKey       = "match"
	finishTypeKey      = "finish"
)

type Pumper struct {
	id          string
	logger      *zap.Logger
	db          *gorm.DB
	rdb         *redis.Client
	lock        *sync.Mutex
	etcdcli     *clientv3.Client
	fetcher     fetcher.Fetcher
	scheduler   scheduler.Scheduler
	entryMap    map[string]map[string]*riotmodel.LeagueEntryDTO // entryMap[Loc][SummonerName]
	sumnMap     map[string]map[string]*riotmodel.SummonerDTO    // sumnMap[Loc][SummonerName]
	matchMap    map[string]map[string]bool                      // matchMap[Loc][matchID]
	Exit        chan struct{}
	out         chan *DBResult
	entrieIdx   []uint
	summonerIdx []uint
	stgy        *Options
}

type DBResult struct {
	Type  string
	Brief string
	Data  interface{}
}

func NewPumper(id string, opts ...Option) (*Pumper, error) {
	stgy := defaultStrategy
	for _, opt := range opts {
		opt(stgy)
	}
	// init
	entryMap := make(map[string]map[string]*riotmodel.LeagueEntryDTO)
	sumnMap := make(map[string]map[string]*riotmodel.SummonerDTO)
	matchMap := make(map[string]map[string]bool)
	for _, l := range stgy.Loc {
		loc, _ := utils.ConvertLocationToLoHoSTR(l)
		entryMap[loc] = make(map[string]*riotmodel.LeagueEntryDTO)
		sumnMap[loc] = make(map[string]*riotmodel.SummonerDTO)
		matchMap[loc] = make(map[string]bool)
	}

	// get deault token
	if stgy.Token == "" {
		workDir := global.ChaConf.DirTree.WorkDir
		filename := "api_key"
		path := filepath.Join(workDir, filename)
		buff, err := os.ReadFile(path)
		if err != nil {
			global.ChaLogger.Error("get api key failed",
				zap.Error(err))
		}
		stgy.Token = string(buff)
	}

	// setup etcd client
	endpoints := []string{stgy.registryURL}
	cli, err := clientv3.New(clientv3.Config{Endpoints: endpoints})
	if err != nil {
		return nil, err
	}

	if global.ChaEnv == global.TestEnv {
		stgy.MaxMatchCount = 1
	}

	return &Pumper{
		id:          id,
		logger:      global.ChaLogger,
		db:          global.ChaDB,
		rdb:         global.ChaRDB,
		lock:        &sync.Mutex{},
		fetcher:     fetcher.NewBrowserFetcher(fetcher.WithToken(stgy.Token)),
		scheduler:   scheduler.NewSchdule(),
		entrieIdx:   make([]uint, 16),
		summonerIdx: make([]uint, 16),
		Exit:        make(chan struct{}),
		out:         make(chan *DBResult),
		entryMap:    entryMap,
		sumnMap:     sumnMap,
		matchMap:    matchMap,
		etcdcli:     cli,

		stgy: stgy,
	}, nil
}

func (p *Pumper) Schedule() {
	p.scheduler.Schedule()
}

func (p *Pumper) handleResult() {
	for result := range p.out {
		switch result.Type {
		case finishTypeKey:
			p.logger.Info(fmt.Sprintf("all %s result store done", result.Brief))
			p.Exit <- struct{}{}
			continue

		case EntryTypeKey:
			entries := result.Data.([]*riotmodel.LeagueEntryDTO)
			if err := p.db.Save(entries).Error; err != nil {
				p.logger.Error("riot entry store failed", zap.Error(err))
			}

		case SummonerTypeKey:
			summoners := result.Data.([]*riotmodel.SummonerDTO)
			if err := p.db.Save(summoners).Error; err != nil {
				p.logger.Error("riot summoner model store failed", zap.Error(err))
			}

		case MatchTypeKey:
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

func (p *Pumper) StartEngine() {
	// go p.LoadAll()
	go p.Schedule()
	// get task from etcd
	go p.getTask()
	go p.watchTasks()

	go p.fetch()
	go p.handleResult()
}

func (p *Pumper) LoadAll() {
	for _, loc := range p.stgy.Loc {
		p.loadSummoners(loc)
		p.loadEntrie(loc)
		p.loadMatch(loc)
	}
}

// UpdateAll 's push task
func (p *Pumper) UpdateAll() {
	p.UpdateEntries()
	p.UpdateSumoner()
	p.UpdateMatch()
	<-p.Exit
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
	endTier, endRank = ConvertRankToStr(p.stgy.TestEndMark1, p.stgy.TestEndMark2)
	// catch panic
	defer func() {
		if err := recover(); err != nil {
			p.logger.Error("fetcher panic",
				zap.Any("err", err),
				zap.String("stack", string(debug.Stack())))
		}
	}()

	for {
		req := p.scheduler.Pull()
		if req.Data == nil {
			// send finish signal
			p.out <- &DBResult{
				Type:  finishTypeKey,
				Brief: req.Type,
				Data:  nil,
			}
			continue
		}

		// fetch and parse
		switch req.Type {
		case bestEntryTypeKey:
			data := req.Data.(*entryTask)
			// api: /lol/league/v4/{BEST}leagues/by-queue/{queue}
			if buff, err = p.fetcher.Get(req.URL); err != nil || buff == nil {
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
				e.Loc = req.Loc
				e.Tier = list.Tier
				e.LeagueID = list.LeagueID
				e.QueType = list.Queue
			}
			// shrink size if test
			if global.ChaEnv == global.TestEnv {
				entries = entries[:testSize]
			}
			// clear list
			list = nil
			p.logger.Info(fmt.Sprintf("all %d %s @ %s data fetch done",
				len(entries), data.Tier, data.Queue[7:]))
			p.handleEntries(entries, req.Loc)
			p.cacheEntries(entries, req.Loc)
			if data.Tier == endTier && data.Rank == endRank {
				p.out <- &DBResult{
					Type:  finishTypeKey,
					Brief: "entry",
					Data:  nil,
				}
				continue
			}

		case mortalEntryTypeKey:
			data := req.Data.(*entryTask)
			for page := 1; ; page++ {
				// api: /lol/league/v4/entries/{queue}/{tier}/{division}
				url := fmt.Sprintf("%s?page=%s", req.URL, strconv.Itoa(page))
				if buff, err = p.fetcher.Get(url); err != nil {
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
					p.logger.Info(fmt.Sprintf("fetch %s %s @ %s,page %d done",
						data.Tier, data.Rank, data.Queue[7:], page))
				}
				for _, e := range entries {
					e.Loc = req.Loc
				}
				// shrink size if test
				if global.ChaEnv == global.TestEnv {
					entries = entries[:testSize]
				}

				p.handleEntries(entries, req.Loc)
				p.cacheEntries(entries, req.Loc)

				// test
				if (global.ChaEnv == global.TestEnv && page == testSize) || len(entries) == 0 {
					p.logger.Info(fmt.Sprintf("all %s %s data fetch done at page %d",
						data.Tier, data.Rank, page))
					break
				}
			}

			if data.Tier == endTier && data.Rank == endRank {
				p.out <- &DBResult{
					Type:  finishTypeKey,
					Brief: "entry",
					Data:  nil,
				}
				continue
			}

		case SummonerTypeKey:
			data := req.Data.(*summonerTask)
			if buff, err = p.fetcher.Get(req.URL); err != nil || buff == nil {
				p.logger.Error(fmt.Sprintf("fetch %s failed",
					data.summonerBrief), zap.Error(err))
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

		case MatchTypeKey:
			data := req.Data.(*matchTask)
			// get old & cur match list
			if buff, err = p.fetcher.Get(req.URL); err != nil || buff == nil {
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
			loc := utils.ConvertLocStrToLocation(req.Loc)
			region := utils.ConvertLocationToRegionHost(loc)
			// fetch match
			for _, matchID := range curMatchList {
				if _, has := p.matchMap[req.Loc][matchID]; has {
					continue
				}
				if _, has := oldMatchList[matchID]; has {
					continue
				} else {
					p.matchMap[req.Loc][matchID] = true
					if tmp := p.FetchMatchByID(req, region, matchID); tmp != nil {
						matches = append(matches, tmp)
					}
				}
			}
			p.logger.Debug(fmt.Sprintf("updating %s's match list @ %d,store %d matches",
				summoner.Name, cnt, len(matches)))
			if cnt == int(p.summonerIdx[loc]) {
				p.rdb.HSet(context.Background(), "lastupdate", "pumper", time.Now().Unix())
			}
			if len(matches) == 0 {
				continue
			}
			p.handleMatches(matches, summoner.Name)
		}
	}
}
