package anres

import (
	"encoding/json"

	"github.com/cralack/ChaosMetrics/server/model"
)

type ChampionBrief struct {
	// basic info
	Image    *model.Image `json:"image"` // 图像
	MetaName string       `json:"id"`    // 英雄ID:Aatrox
	// statistical data
	WinRate        float32 `json:"win_rate"`         // 胜率 30%
	PickRate       float32 `json:"pick_rate"`        // 登场率 15%
	BanRate        float32 `json:"ban_rate"`         // Ban率
	AvgDamageDealt float32 `json:"avg_damage_dealt"` // 场均输出占比 12%
	AvgDeadTime    float32 `json:"avg_dead_time"`    // 场均死亡时长
}

func (c *ChampionBrief) MarshalBinary() ([]byte, error) {
	return json.Marshal(c)
}
func (c *ChampionBrief) UnmarshalBinary(bt []byte) error {
	return json.Unmarshal(bt, c)
}
