package analyzer

import (
	"github.com/cralack/ChaosMetrics/server/model/anres"
)

type Judger func(tar *anres.ChampionDetail)

func JudgeARAM() Judger {
	return func(tar *anres.ChampionDetail) {
		tar.RankScore = tar.WinRate*1000 +
			tar.PickRate*400 +
			tar.AvgKDA*5 +
			tar.AvgKP*50 +
			tar.AvgDamageDealt/1000*1.2 +
			tar.AvgDamageTaken/1000 +
			tar.AvgTimeCCing*0.8
	}
}

func JudgeClassic() Judger {
	return func(tar *anres.ChampionDetail) {
		tar.RankScore = tar.WinRate*1000 +
			tar.PickRate*400 +
			tar.AvgKDA*5 +
			tar.AvgKP*50 +
			tar.AvgDamageDealt/1000*1.2 +
			tar.AvgDamageTaken/1000 +
			tar.AvgTimeCCing*0.5 +
			tar.AvgVisionScore*0.3
	}
}
