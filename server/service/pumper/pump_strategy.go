package pumper

import (
	"github.com/cralack/ChaosMetrics/server/model/riotmodel"
)

type RiotStrategy struct {
	Loc     []uint
	Que     []uint
	Tier    []uint // todo
	MaxSize int
}
type Option func(stgy *RiotStrategy)

var defaultStrategy = &RiotStrategy{
	Loc:     []uint{riotmodel.TW2},
	Que:     []uint{riotmodel.RANKED_SOLO_5x5},
	MaxSize: 400,
}

func WithLocs(loc []uint) Option {
	return func(stgy *RiotStrategy) {
		stgy.Loc = loc
	}
}

func WithQues(que []uint) Option {
	return func(stgy *RiotStrategy) {
		stgy.Que = que
	}
}

func WithTier(tier []uint) Option {
	return func(stgy *RiotStrategy) {
		stgy.Tier = tier
	}
}
