package riotmodel

import (
	"encoding/json"
	"time"
)

type ChampionMasteryDto struct {
	ChampionPointsUntilNextLevel int64     `json:"championPointsUntilNextLevel"` // 距离下一等级所需的点数。如果玩家达到了该英雄的最高等级，则为零。
	ChestGranted                 bool      `json:"chestGranted"`                 // 在当前赛季是否已经获得了该英雄的宝箱。
	ChampionID                   int64     `json:"championId"`                   // 该记录对应的英雄ID。
	LastPlayTime                 time.Time `json:"lastPlayTime"`                 // 玩家最后一次玩该英雄的时间，以Unix毫秒时间格式表示。
	ChampionLevel                int       `json:"championLevel"`                // 指定玩家和英雄组合的英雄等级。
	SummonerID                   string    `json:"summonerId"`                   // 该记录对应的召唤师ID（加密）。
	ChampionPoints               int       `json:"championPoints"`               // 玩家和英雄组合的总英雄点数，用于确定英雄等级。
	ChampionPointsSinceLastLevel int64     `json:"championPointsSinceLastLevel"` // 达到当前等级后获得的点数。
	TokensEarned                 int       `json:"tokensEarned"`                 // 当前英雄等级下获得的代币数。当英雄等级提升时，tokensEarned重置为0。
}

var _ DTO = &ChampionMasteryDto{}

func (p *ChampionMasteryDto) UnmarshalJSON(data []byte) error {
	var f map[string]interface{}
	err := json.Unmarshal(data, &f)
	if err != nil {
		return err
	}
	for k, v := range f {
		switch k {
		case "championPointsUntilNextLevel":
			p.ChampionPointsUntilNextLevel = int64(v.(float64))
		case "chestGranted":
			p.ChestGranted = v.(bool)
		case "championId":
			p.ChampionID = int64(v.(float64))
		case "lastPlayTime":
			// Assuming lastPlayTime is in milliseconds
			lastPlayTimeMillis := int64(v.(float64))
			p.LastPlayTime = time.Unix(0, lastPlayTimeMillis*int64(time.Millisecond)).UTC()
		case "championLevel":
			p.ChampionLevel = int(v.(float64))
		case "summonerId":
			p.SummonerID = v.(string)
		case "championPoints":
			p.ChampionPoints = int(v.(float64))
		case "championPointsSinceLastLevel":
			p.ChampionPointsSinceLastLevel = int64(v.(float64))
		case "tokensEarned":
			p.TokensEarned = int(v.(float64))
		}
	}

	return nil
}
