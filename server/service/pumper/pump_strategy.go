package pumper

import (
	"fmt"
	"time"
	
	"github.com/cralack/ChaosMetrics/server/model/riotmodel"
)

type RiotStrategy struct {
	Loc []uint
	Que []uint
	EndMark  []uint
	MaxSize  int
	LifeTime time.Duration
}

type Option func(stgy *RiotStrategy)

var defaultStrategy = &RiotStrategy{
	Loc:      []uint{riotmodel.TW2},
	Que:      []uint{riotmodel.RANKED_SOLO_5x5},
	EndMark:  []uint{riotmodel.IRON, 4},
	MaxSize:  500,
	LifeTime: time.Hour * 12,
	// LifeTime: -1, // cache forever
}

//	Example:WithLoc(riotmodel.BR1,riotmodel.EUN1)
func WithLoc(locs ...uint) Option {
	return func(stgy *RiotStrategy) {
		stgy.Loc = make([]uint, 0, 16)
		for _, loc := range locs {
			stgy.Loc = append(stgy.Loc, loc)
		}
	}
}

//	Example:WithAreaLoc(riotmodel.LOC_ALL)
//	Example:WithAreaLoc(riotmodel.LOC_AMERICAS,riotmodel.LOC_ASIA)
func WithAreaLoc(locs ...uint) Option {
	return func(stgy *RiotStrategy) {
		res := make([]uint, 0, 16)
		america := []uint{
			riotmodel.BR1,
			riotmodel.LA1,
			riotmodel.LA2,
			riotmodel.NA1,
		}
		asia := []uint{
			riotmodel.KR1,
			riotmodel.JP1,
		}
		europe := []uint{
			riotmodel.EUN1,
			riotmodel.EUW1,
			riotmodel.TR1,
			riotmodel.RU,
		}
		sea := []uint{
			riotmodel.OC1,
			riotmodel.PH2,
			riotmodel.SG2,
			riotmodel.TH2,
			riotmodel.TW2,
			riotmodel.VN2,
		}
		
		for loc := range locs {
			switch loc {
			case riotmodel.LOC_AMERICAS:
				res = append(res, america...)
			case riotmodel.LOC_ASIA:
				res = append(res, asia...)
			case riotmodel.LOC_EUROPE:
				res = append(res, europe...)
			case riotmodel.LOC_SEA:
				res = append(res, sea...)
			case riotmodel.LOC_ALL:
				stgy.Loc = make([]uint, 0, 16)
				stgy.Loc = append(stgy.Loc, america...)
				stgy.Loc = append(stgy.Loc, asia...)
				stgy.Loc = append(stgy.Loc, europe...)
				stgy.Loc = append(stgy.Loc, sea...)
				return
			default:
				panic(fmt.Sprintf("unknown location option: %d", loc))
			}
		}
		stgy.Loc = res
	}
}

//	Example:WithQues(riotmodel.RANKED_SOLO_5x5)
func WithQues(que []uint) Option {
	return func(stgy *RiotStrategy) {
		stgy.Que = que
	}
}

//	Example:WithEndMark(riotmodel.DIAMOND,1)
//	Example:WithEndMark(riotmodel.IRON,4)
func WithEndMark(tier, div uint) Option {
	return func(stgy *RiotStrategy) {
		stgy.EndMark = []uint{tier, div}
	}
}