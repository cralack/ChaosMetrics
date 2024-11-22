package updater

import (
	"time"

	"github.com/cralack/ChaosMetrics/server/internal/global"
	"github.com/cralack/ChaosMetrics/server/model/riotmodel"
)

type Strategy struct {
	Loc         []riotmodel.LOCATION
	Lang        []riotmodel.LANG
	LifeTime    time.Duration
	ForceUpdate bool
	EndMark     string
}

type Setup func(stgy *Strategy) // Strategy的配置选项

var defaultStrategy = &Strategy{
	Loc:  []riotmodel.LOCATION{riotmodel.TW2},
	Lang: []riotmodel.LANG{riotmodel.LANG_zh_CN, riotmodel.LANG_en_US},
	// Lang:     []uint{riotmodel.LANG_zh_CN},
	LifeTime: time.Hour * 24 * 7, // 7 day
}

func _(locs ...riotmodel.LOCATION) Setup {
	return func(stgy *Strategy) {
		tmp := make([]riotmodel.LOCATION, 0, 16)
		for _, loc := range locs {
			if 16 < loc {
				global.ChaLogger.Error("wrong param,loc need < 16,using default option")
				return
			}
			tmp = append(tmp, loc)
		}
		stgy.Loc = tmp
	}
}

func WithLifeTime(life time.Duration) Setup {
	return func(stgy *Strategy) {
		stgy.LifeTime = life
	}
}

func WithForceUpdate(flag bool) Setup {
	return func(stgy *Strategy) {
		stgy.ForceUpdate = flag
	}
}

func WithEndmark(version string) Setup {
	return func(stgy *Strategy) {
		stgy.EndMark = version
	}
}
