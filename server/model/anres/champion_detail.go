package anres

import (
	"encoding/json"

	"github.com/cralack/ChaosMetrics/server/model"
)

type ChampionDetail struct {
	// basic info
	ID       string       `json:"idx"`       // 存储key
	Image    *model.Image `json:"image"`     // 图像
	Loc      string       `json:"loc"`       // 对局服务器
	Version  string       `json:"version"`   // 对局版本
	MetaName string       `json:"id"`        // 英雄ID:Aatrox
	Key      string       `json:"key"`       // 英雄Key:266
	Name     string       `json:"name"`      // 英雄名称:暗裔剑魔
	Title    string       `json:"title"`     // 英雄称号:亚托克斯
	GameMode string       `json:"game_mode"` // 游戏模式
	// statistical data
	RankScore      float32 `json:"rank_score"`       // 英雄综合打分
	TotalWin       float32 `json:"total_win"`        // 英雄胜利总数
	TotalPlayed    float32 `json:"total_played"`     // 英雄出场总数
	WinRate        float32 `json:"win_rate"`         // 胜率 30%
	PickRate       float32 `json:"pick_rate"`        // 登场率 15%
	BanRate        float32 `json:"ban_rate"`         // Ban率
	AvgKDA         float32 `json:"avg_kda"`          // 场均KDA 15%
	AvgKP          float32 `json:"avg_kp"`           // 场均参团率 10%
	AvgDamageDealt float32 `json:"avg_damage_dealt"` // 场均输出占比 12%
	AvgDamageTaken float32 `json:"avg_damage_taken"` // 场均承伤占比 10%
	AvgTimeCCing   float32 `json:"avg_time_ccing"`   // 场均控制时长 5%
	AvgVisionScore float32 `json:"avg_vision_score"` // 场均视野得分 3%
	AvgDeadTime    float32 `json:"avg_dead_time"`    // 场均死亡时长
	// build winrate
	ItemWin  map[string]map[string]int `json:"item"`
	PerkWin  map[string]int            `json:"perk"`
	SkillWin map[string]int            `json:"skill"`
	SpellWin map[string]int            `json:"spell"`
}

func (c *ChampionDetail) TableName() string {
	return "analyzed_champions"
}
func (c *ChampionDetail) MarshalBinary() ([]byte, error) {
	return json.Marshal(c)
}

func (c *ChampionDetail) UnmarshalBinary(bt []byte) error {
	return json.Unmarshal(bt, c)
}
