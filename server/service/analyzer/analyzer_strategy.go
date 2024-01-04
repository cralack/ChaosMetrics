package analyzer

import (
	"github.com/cralack/ChaosMetrics/server/internal/global"
	"github.com/cralack/ChaosMetrics/server/model/riotmodel"
)

type options struct {
	Loc       []uint // 地区列表
	Mode      []uint // 游戏模式
	Lang      []uint // 语言
	BatchSize int
}

type Option func(stgy *options) // Strategy的配置选项

var defaultOptions = &options{
	Loc:       []uint{riotmodel.TW2},                              // 默认地区为台湾
	Mode:      []uint{riotmodel.ARAM},                             // 默认模式为大乱斗
	Lang:      []uint{riotmodel.LANG_zh_CN, riotmodel.LANG_en_US}, // 默认中文、英文
	BatchSize: 100,
	// LifeTime: -1, // cache forever
}

// Example:WithLoc(riotmodel.BR1,riotmodel.EUN1)
func WithLoc(locs ...uint) Option {
	return func(stgy *options) {
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

func WithMode(mode ...uint) Option {
	return func(stgy *options) {
		tmp := make([]uint, 0, 4)
		for _, m := range mode {
			if riotmodel.CHERRY < m {
				global.GVA_LOG.Error("wrong param,using default option")
				return
			}
			tmp = append(tmp, m)
		}
		stgy.Mode = tmp
	}
}
