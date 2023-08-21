package pumper

import (
	"fmt"
	"sync"
	
	"github.com/cralack/ChaosMetrics/server/global"
	"github.com/cralack/ChaosMetrics/server/model/riotmodel"
	"github.com/cralack/ChaosMetrics/server/service/fetcher"
	"github.com/cralack/ChaosMetrics/server/service/scheduler"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

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
	out         chan *ParseResult
	entrieIdx   []uint
	summonerIdx []uint
	stgy        *Strategy
}

type ParseResult struct {
	Type  string
	Brief string
	Data  interface{}
}

func NewPumper(opts ...Option) *Pumper {
	stgy := defaultStrategy
	for _, opt := range opts {
		opt(stgy)
	}
	
	return &Pumper{
		logger: global.GVA_LOG,
		db:     global.GVA_DB,
		rdb:    global.GVA_RDB,
		lock:   &sync.Mutex{},
		fetcher: fetcher.NewBrowserFetcher(
			fetcher.WithAPIToken(stgy.Token),
		),
		scheduler: scheduler.NewSchdule(),
		
		entrieIdx:   make([]uint, 16),
		summonerIdx: make([]uint, 16),
		out:         make(chan *ParseResult),
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

func (p *Pumper) UpdateAll() {
	exit := make(chan struct{})
	go p.Schedule()
	go p.handleResult(exit)
	
	p.UpdateEntries(exit)
	p.UpdateSumoner(exit)
	p.UpdateMatch(exit)
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
