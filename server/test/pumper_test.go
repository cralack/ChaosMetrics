package test

import (
	"testing"
	"time"
	
	"github.com/cralack/ChaosMetrics/server/global"
	"github.com/cralack/ChaosMetrics/server/model/riotmodel"
	"github.com/cralack/ChaosMetrics/server/service/pumper"
	"go.uber.org/zap"
)

func Test_pumper(t *testing.T) {
	puper := pumper.NewPumper()
	go puper.UpdateEntries()
	time.Sleep(time.Second * 60)
	Test_store_check(t)
}

func Test_store_check(t *testing.T) {
	var tier uint
	var count int64
	logger := global.GVA_LOG
	db := global.GVA_DB
	res := db.Model(&riotmodel.LeagueEntryDTO{}).Count(&count)
	if err := res.Error; err != nil {
		logger.Error("get idx failed", zap.Error(err))
	} else {
		logger.Info("total count", zap.Int64("count", count))
	}
	for tier = riotmodel.CHALLENGER; tier <= riotmodel.IRON; tier++ {
		tierSTR, _ := pumper.ConvertRankToStr(tier, 1)
		
		var cnt int64
		res := global.GVA_DB.Model(&riotmodel.LeagueEntryDTO{}).Where("tier=?", tierSTR).Count(&cnt)
		if res.Error != nil {
			t.Log(res.Error)
		} else {
			t.Log(tierSTR, cnt)
		}
	}
}
