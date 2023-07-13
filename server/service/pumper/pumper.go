package pumper

import (
	"sync"
	"time"
	
	"github.com/cralack/ChaosMetrics/server/global"
	"github.com/cralack/ChaosMetrics/server/model/riotmodel"
	"github.com/cralack/ChaosMetrics/server/service/fetcher"
	"github.com/cralack/ChaosMetrics/server/service/scheduler"
	"github.com/cralack/ChaosMetrics/server/utils"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Pumper struct {
	logger     *zap.Logger
	db         *gorm.DB
	lock       *sync.Mutex
	rdb        *redis.Client
	fetcher    fetcher.Fetcher
	scheduler  scheduler.Scheduler
	entriesDic map[string]map[string]*riotmodel.LeagueEntryDTO // entriesDic[Loc][SummonerID]
	rater      chan struct{}
	out        chan *ParseResult
	entrieIdx  []uint
	stgy       *RiotStrategy
}

type ParseResult struct {
	Type  string
	Brief string
	// Idx   uint
	Page  int
	Data  interface{}
}

func NewPumper(opts ...Option) *Pumper {
	stgy := defaultStrategy
	for _, opt := range opts {
		opt(stgy)
	}
	logger := global.GVA_LOG
	db := global.GVA_DB
	idx := make([]uint, 16)
	// calculate each loc idx
	var count int64
	for _, loc := range stgy.Loc {
		locStr, _ := utils.ConvertHostURL(loc)
		res := db.Model(&riotmodel.LeagueEntryDTO{Loc: locStr}).Count(&count)
		if err := res.Error; err != nil {
			logger.Error("get idx failed", zap.Error(err))
		}
		idx[loc] = uint(count)
	}
	
	return &Pumper{
		logger:     logger,
		db:         db,
		rdb:        global.GVA_RDB,
		lock:       &sync.Mutex{},
		entrieIdx:  idx,
		out:        make(chan *ParseResult),
		rater:      make(chan struct{}),
		entriesDic: make(map[string]map[string]*riotmodel.LeagueEntryDTO),
		fetcher:    fetcher.NewBrowserFetcher(),
		scheduler:  scheduler.NewSchdule(),
		stgy:       stgy,
	}
}
func (p *Pumper) Schedule() {
	p.scheduler.Schedule()
}

func (p *Pumper) startTimer() {
	ticker := time.NewTicker(time.Millisecond * 500)
	for {
		<-ticker.C
		if p.fetcher.TryAcquire() {
			p.rater <- struct{}{}
		}
	}
}

func (p *Pumper) UpdateAll() {
	p.InitEntries()
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
