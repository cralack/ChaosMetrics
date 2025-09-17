package test

import (
	"context"
	"testing"

	"github.com/cralack/ChaosMetrics/server/internal/service/pumper"
	"github.com/cralack/ChaosMetrics/server/model/riotmodel"
)

func Test_pumper_update(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	p, _ := pumper.NewPumper(
		"1",
		pumper.WithLoc(riotmodel.NA1),
		pumper.WithEndMark(riotmodel.DIAMOND, 1),
		pumper.WithQues(riotmodel.RANKED_SOLO_5x5),
		pumper.WithMaxMatchCount(3),
		pumper.WithContext(ctx),
	)
	p.StartEngine()
	p.UpdateAll()
}
