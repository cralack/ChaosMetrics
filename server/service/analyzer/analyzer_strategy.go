package analyzer

import (
	"github.com/cralack/ChaosMetrics/server/internal/global"
	"github.com/cralack/ChaosMetrics/server/model/riotmodel"
)

type options struct {
	Loc       []riotmodel.LOCATION // 地区列表
	Mode      []riotmodel.GAMEMODE // 游戏模式
	Lang      []riotmodel.LANG     // 语言
	BatchSize int
}

type Option func(stgy *options) // Strategy的配置选项

var defaultOptions = &options{
	Loc:       []riotmodel.LOCATION{riotmodel.TW2},                          // 默认地区为台湾
	Mode:      []riotmodel.GAMEMODE{riotmodel.ARAM},                         // 默认模式为大乱斗
	Lang:      []riotmodel.LANG{riotmodel.LANG_zh_CN, riotmodel.LANG_en_US}, // 默认中文、英文
	BatchSize: 100,
	// LifeTime: -1, // cache forever
}

func WithLoc(locs ...riotmodel.LOCATION) Option {
	return func(stgy *options) {
		tmp := make([]riotmodel.LOCATION, 0, 16)
		for _, loc := range locs {
			if 16 < loc {
				global.GvaLog.Error("wrong param,loc need < 16,using default option")
				return
			}
			tmp = append(tmp, loc)
		}
		stgy.Loc = tmp
	}
}

func WithMode(mode ...riotmodel.GAMEMODE) Option {
	return func(stgy *options) {
		tmp := make([]riotmodel.GAMEMODE, 0, 4)
		for _, m := range mode {
			if riotmodel.CHERRY < m {
				global.GvaLog.Error("wrong param,using default option")
				return
			}
			tmp = append(tmp, m)
		}
		stgy.Mode = tmp
	}
}
