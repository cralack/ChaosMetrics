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
	// Test_db_store_check(t)
	// Test_check(t)
}
