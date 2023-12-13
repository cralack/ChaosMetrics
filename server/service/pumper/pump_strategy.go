package pumper

import (
	"fmt"
	"time"

	"github.com/cralack/ChaosMetrics/server/global"
	"github.com/cralack/ChaosMetrics/server/model/riotmodel"
)

type Strategy struct {
	Token         string
	Loc           []uint        // 地区列表
	Que           []uint        // 队列类型列表
	TestEndMark   []uint        // 测试用终止标记
	MaxSize       int           // Task最大切割尺寸
	MaxMatchCount int           // 最大比赛场次
	Retry         uint          // 任务重试次数
	LifeTime      time.Duration // 缓存生命周期
}

type Option func(stgy *Strategy) // Strategy的配置选项

var defaultStrategy = &Strategy{
	Token:         "",
	Loc:           []uint{riotmodel.TW2},             // 默认地区为台湾
	Que:           []uint{riotmodel.RANKED_SOLO_5x5}, // 默认队列类型为排位赛5v5
	TestEndMark:   []uint{riotmodel.DIAMOND, 1},      // 默认终止标记为钻I
	MaxSize:       500,                               // 默认任务切割尺寸为500
	MaxMatchCount: 20,                                // 默认读取最近20场比赛
	Retry:         3,                                 // 默认单个任务重试次数3
	LifeTime:      time.Hour * 24,                    // 默认缓存生命周期为24小时
	// LifeTime: -1, // cache forever
}

// Example:WithLoc(riotmodel.BR1,riotmodel.EUN1)
func WithLoc(locs ...uint) Option {
	return func(stgy *Strategy) {
		tmp := make([]uint, 0, 16)
		for _, loc := range locs {
			if 16 < loc {
				global.GVA_LOG.Error("wrong param,loc need < 16,using default option")
				return
			}
			tmp = append(tmp, loc)
		}
		stgy.Loc = tmp
	}
}

// Example:WithAreaLoc(riotmodel.LOC_ALL)
// Example:WithAreaLoc(riotmodel.LOC_AMERICAS,riotmodel.LOC_ASIA)
func WithAreaLoc(areas ...uint) Option {
	return func(stgy *Strategy) {
		tmp := make([]uint, 0, 16)
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

		for _, area := range areas {
			if 4 < area {
				global.GVA_LOG.Error("wrong param,area need < 4,using default option")
				return
			}
			switch area {
			case riotmodel.LOC_AMERICAS:
				tmp = append(tmp, america...)
			case riotmodel.LOC_ASIA:
				tmp = append(tmp, asia...)
			case riotmodel.LOC_EUROPE:
				tmp = append(tmp, europe...)
			case riotmodel.LOC_SEA:
				tmp = append(tmp, sea...)
			case riotmodel.LOC_ALL:
				stgy.Loc = make([]uint, 0, 16)
				stgy.Loc = append(stgy.Loc, america...)
				stgy.Loc = append(stgy.Loc, asia...)
				stgy.Loc = append(stgy.Loc, europe...)
				stgy.Loc = append(stgy.Loc, sea...)
				return
			default:
				panic(fmt.Sprintf("unknown location option: %d", area))
			}
		}
		stgy.Loc = tmp
	}
}

// Example:WithQues(riotmodel.RANKED_SOLO_5x5)
func WithQues(ques ...uint) Option {
	return func(stgy *Strategy) {
		tmp := make([]uint, 0, 3)
		for _, que := range ques {
			if 3 < que {
				global.GVA_LOG.Error("wrong param,que need < 3,using default option")
				return
			}
			tmp = append(tmp, que)
		}
		stgy.Que = tmp
	}
}

// Example:WithEndMark(riotmodel.DIAMOND,1)
// Example:WithEndMark(riotmodel.IRON,4)
func WithEndMark(tier, div uint) Option {
	return func(stgy *Strategy) {
		if riotmodel.IRON < tier || 4 < div || div < 1 {
			global.GVA_LOG.Error("wrong param,end mark need DIAMON <= tier <= IRON" +
				" && I <= div <= IV.using default option")
			return
		}
		stgy.TestEndMark = []uint{tier, div}
	}
}

func WithToken(token string) Option {
	return func(stgy *Strategy) {
		if token != "" {
			stgy.Token = token
		}
	}
}
