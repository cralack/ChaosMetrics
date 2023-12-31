package updater

import (
	"time"
	
	"github.com/cralack/ChaosMetrics/server/global"
	"github.com/cralack/ChaosMetrics/server/model/riotmodel"
)

type Strategy struct {
	Loc      []uint
	Lang     []uint
	LifeTime time.Duration
}

type Option func(stgy *Strategy) // Strategy的配置选项

var defaultStrategy = &Strategy{
	Loc:  []uint{riotmodel.TW2},
	Lang: []uint{riotmodel.LANG_zh_CN, riotmodel.LANG_en_US},
	// Lang:     []uint{riotmodel.LANG_zh_CN},
	LifeTime: time.Hour * 24 * 7, // 7 day
}

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

func WithLifeTime(life time.Duration) Option {
	return func(stgy *Strategy) {
		stgy.LifeTime = life
	}
}
