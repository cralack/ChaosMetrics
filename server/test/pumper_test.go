package test

import (
	"testing"
	
	"github.com/cralack/ChaosMetrics/server/model/riotmodel"
	"github.com/cralack/ChaosMetrics/server/service/pumper"
)

func Test_pumper_update(t *testing.T) {
	p := pumper.NewPumper(
		"",
		pumper.WithEndMark(riotmodel.DIAMOND, 4),
	)
	
	p.UpdateAll()
	
}