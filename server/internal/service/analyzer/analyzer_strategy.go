package analyzer

import (
	"github.com/cralack/ChaosMetrics/server/internal/global"
	"github.com/cralack/ChaosMetrics/server/model/riotmodel"
	"github.com/cralack/ChaosMetrics/server/utils"
)

type options struct {
	Loc       []riotmodel.LOCATION // 地区列表
	Mode      []riotmodel.GAMEMODE // 游戏模式
	Lang      []riotmodel.LANG     // 语言
	Versions  []string
	BatchSize int
}

type Option func(stgy *options) // Strategy的配置选项

var defaultOptions = &options{
	Loc:       []riotmodel.LOCATION{riotmodel.TW2},                          // 默认地区为台湾
	Mode:      []riotmodel.GAMEMODE{riotmodel.ARAM, riotmodel.CLASSIC},      // 默认模式为大乱斗
	Lang:      []riotmodel.LANG{riotmodel.LANG_zh_CN, riotmodel.LANG_en_US}, // 默认中文、英文
	BatchSize: 1000,
	Versions:  make([]string, 0),
	// LifeTime: -1, // cache forever
}

func WithLoc(locs ...riotmodel.LOCATION) Option {
	return func(stgy *options) {
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

func WithMode(mode ...riotmodel.GAMEMODE) Option {
	return func(stgy *options) {
		tmp := make([]riotmodel.GAMEMODE, 0, 10)
		for _, m := range mode {
			if riotmodel.CHERRY < m {
				global.ChaLogger.Error("wrong param,using default option")
				return
			}
			tmp = append(tmp, m)
		}
		stgy.Mode = tmp
	}
}

func WithVersions(versions ...string) Option {
	return func(stgy *options) {
		for _, ver := range versions {
			if ver == "all" {
				stgy.Versions = utils.GetCurMajorVersions()
				return
			}
		}
		stgy.Versions = versions
	}
}
