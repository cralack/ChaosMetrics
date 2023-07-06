package pumper

import (
	"encoding/json"
	"fmt"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	
	"github.com/cralack/ChaosMetrics/server/global"
	"github.com/cralack/ChaosMetrics/server/model/riotmodel"
	"github.com/cralack/ChaosMetrics/server/service/fetcher"
	"github.com/cralack/ChaosMetrics/server/service/scheduler"
	"github.com/cralack/ChaosMetrics/server/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Pumper struct {
	logger    *zap.Logger
	db        *gorm.DB
	lock      *sync.Mutex
	fetcher   fetcher.Fetcher
	scheduler scheduler.Scheduler
	out       chan *ParseResult
	entrieIdx uint
	stgy      *RiotStrategy
}

type ParseResult struct {
	Type  string
	Brief string
	Idx   uint
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
	var count int64
	res := db.Model(&riotmodel.LeagueEntryDTO{}).Count(&count)
	if err := res.Error; err != nil {
		logger.Error("get idx failed", zap.Error(err))
	}
	
	return &Pumper{
		logger:    logger,
		db:        db,
		lock:      &sync.Mutex{},
		entrieIdx: uint(count) + 1,
		out:       make(chan *ParseResult),
		fetcher:   fetcher.NewBrowserFetcher(),
		scheduler: scheduler.NewSchdule(),
		stgy:      stgy,
	}
}

func (p *Pumper) schedule() {
	p.scheduler.Schedule()
}

func (p *Pumper) UpdateEntries() {
	go p.schedule()
	for _, loc := range p.stgy.Loc {
		for _, que := range p.stgy.Que {
			// Generate URLs
			go p.createEntriesWork(loc, que)
		}
	}
	// go Watch_ETCD
	go p.fetchDTO()
	// DTO store
	p.handleResult()
}

func (p *Pumper) createEntriesWork(loc, que uint) {
	var (
		stem string
		tier uint
		div  uint
	)
	prefix := utils.ConvertPlatformURL(loc)
	queStr := getQueueString(que)
	// generate BEST URL task
	for tier = riotmodel.CHALLENGER; tier <= riotmodel.MASTER; tier++ {
		t, _ := ConvertRankToStr(tier, 1)
		stem = fmt.Sprintf("/lol/league/v4/%sleagues/by-queue/%s",
			strings.ToLower(t), queStr)
		p.scheduler.Push(&scheduler.Task{
			Key:   "best",
			Brief: t,
			URL:   prefix + stem,
		})
	}
	// generate MORTAL URL task
	for tier = riotmodel.DIAMOND; tier <= riotmodel.IRON; tier++ {
		for div = 1; div <= 4; div++ {
			tierStr, divStr := ConvertRankToStr(tier, div)
			stem = fmt.Sprintf("/lol/league/v4/entries/%s/%s/%s",
				queStr, tierStr, divStr)
			p.scheduler.Push(&scheduler.Task{
				Key:   "mortal",
				Brief: tierStr + " " + divStr,
				URL:   prefix + stem,
			})
		}
	}
}

func (p *Pumper) fetchDTO() {
	var (
		page    int
		buff    []byte
		err     error
		list    riotmodel.LeagueListDTO
		entries []*riotmodel.LeagueEntryDTO
	)
	// catch panic
	defer func() {
		if err := recover(); err != nil {
			p.logger.Error("fetcher panic",
				zap.Any("err", err),
				zap.String("stack", string(debug.Stack())))
		}
	}()
	
	for {
		if p.fetcher.TryAcquire() {
			req := p.scheduler.Pull()
			// p.StoreVisited(req.Key)
			p.logger.Info(fmt.Sprintf("fetching request URL:%s", req.URL))
			// fetch and parse
			switch req.Key {
			case "best":
				// api: /lol/league/v4/{BEST}leagues/by-queue/{queue}
				if buff, err = p.fetcher.Get(req.URL); err != nil {
					p.logger.Error(fmt.Sprintf("fetch %s failed", req.Brief),
						zap.Error(err))
				}
				if err = json.Unmarshal(buff, &list); err != nil {
					p.logger.Error(fmt.Sprintf("unmarshal json to %s failed",
						"LeagueListDTO"), zap.Error(err))
				}
				entries = list.Entries
				p.handleEntry(entries, req.Brief)
				// check oversize && chunk
				var chunks [][]*riotmodel.LeagueEntryDTO
				if len(entries) > p.stgy.MaxSize {
					totalSize := len(entries)
					chunkSize := p.stgy.MaxSize
					for i := 0; i < totalSize; i += chunkSize {
						end := i + chunkSize
						if end > totalSize {
							end = totalSize
						}
						chunks = append(chunks, entries[i:end])
					}
				} else {
					chunks = append(chunks, entries)
				}
				for i, chunk := range chunks {
					p.out <- &ParseResult{
						Type:  "entry",
						Brief: req.Brief,
						Page:  i + 1,
						Data:  chunk,
					}
				}
			
			case "mortal":
				page = 0
				for {
					page++
					// api: /lol/league/v4/entries/{queue}/{tier}/{division}
					if buff, err = p.fetcher.Get(fmt.Sprintf("%s?page=%s",
						req.URL, strconv.Itoa(page))); err != nil {
						p.logger.Error(fmt.Sprintf("fetch %s failed", req.Brief),
							zap.Error(err))
					}
					if err = json.Unmarshal(buff, &entries); err != nil {
						p.logger.Error(fmt.Sprintf("unmarshal json to %s failed",
							"LeagueItemDTO"), zap.Error(err))
					}
					if len(entries) == 0 {
						p.logger.Info(fmt.Sprintf("all %s data fetch done at page %02d",
							req.Brief, page))
						break
					}
					tier := strings.Split(req.Brief, " ")[0]
					p.handleEntry(entries, tier)
					p.out <- &ParseResult{
						Type:  "entry",
						Brief: req.Brief,
						Page:  page,
						Data:  entries,
					}
				}
			}
		} else {
			continue
		}
	}
}

func (p *Pumper) handleResult() {
	for result := range p.out {
		switch result.Type {
		case "entry":
			data := result.Data.([]*riotmodel.LeagueEntryDTO)
			size := len(data)
			if size == 0 {
				continue
			}
			// check duplicate
			
			// store data
			// time.Sleep(time.Millisecond * 300)
			if err := p.db.Create(data).Error; err != nil {
				p.logger.Error("riot model store failed", zap.Error(err))
			} else {
				p.logger.Info(fmt.Sprintf("%s's entries store succeed %d", data[0].Tier+" "+data[0].Rank, size))
			}
		}
		
	}
}
func (p *Pumper) handleEntry(entries []*riotmodel.LeagueEntryDTO, tier string) {
	for i, entry := range entries {
		entry.Tier = tier
		entry.ID = p.entrieIdx + uint(i)
	}
	
	p.lock.Lock()
	defer p.lock.Unlock()
	p.entrieIdx += uint(len(entries))
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
