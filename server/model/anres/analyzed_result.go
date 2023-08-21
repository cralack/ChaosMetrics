package anres

import (
	"encoding/json"
	
	"github.com/cralack/ChaosMetrics/server/model"
	"gorm.io/gorm"
)

type Champion struct {
	gorm.Model
	// basic info
	Image    *model.Image `json:"image" gorm:"embedded;embeddedPrefix:icon_image_"` // 图像
	Loc      string       `json:"loc" gorm:"column:loc;type:varchar(20)"`           // 对局服务器
	Version  string       `json:"version" gorm:"version"`                           // 对局版本
	MetaName string       `json:"id" gorm:"column:meta_name;index"`                 // 英雄ID:Aatrox
	Key      string       `json:"key" gorm:"column:key"`                            // 英雄Key:266
	Name     string       `json:"name" gorm:"column:name"`                          // 英雄名称:暗裔剑魔
	Title    string       `json:"title" gorm:"column:title"`                        // 英雄称号:亚托克斯
	GameMode string       `json:"game_mode" gorm:"column:game_mode"`                // 游戏模式
	// statistical data
	RankScore      float32 `json:"rank_score" gorm:"column:rank_score"`             // 英雄综合打分
	TotalWin       float32 `json:"total_win" gorm:"column:total_win"`               // 英雄胜利总数
	TotalPlayed    float32 `json:"total_played" gorm:"column:total_played"`         // 英雄出场总数
	WinRate        float32 `json:"win_rate" gorm:"column:win_rate"`                 // 胜率 30%
	PickRate       float32 `json:"pick_rate" gorm:"column:pick_rate"`               // 登场率 15%
	BanRate        float32 `json:"ban_rate" gorm:"column:ban_rate"`                 // Ban率
	AvgKDA         float32 `json:"avg_kda" gorm:"column:avg_kda"`                   // 场均KDA 15%
	AvgKP          float32 `json:"avg_kp" gorm:"column:avg_kp"`                     // 场均参团率 10%
	AvgDamageDealt float32 `json:"avg_damage_dealt" gorm:"column:avg_damage_dealt"` // 场均输出占比 12%
	AvgDamageTaken float32 `json:"avg_damage_taken" gorm:"column:avg_damage_taken"` // 场均承伤占比 10%
	AvgTimeCCing   float32 `json:"avg_time_ccing" gorm:"column:avg_time_ccing"`     // 场均控制时长 5%
	AvgVisionScore float32 `json:"avg_vision_score" gorm:"column:avg_vision_score"` // 场均视野得分 3%
	AvgDeadTime    float32 `json:"avg_dead_time" gorm:"column:avg_dead_time"`       // 场均死亡时长
	// build winrate
	ItemWin  map[string]map[string]int `json:"item" gorm:"-"`
	ItemSTR  string                    `json:"-" gorm:"item"`
	PerkWin  map[string]int            `json:"perk" gorm:"-"`
	PerkSTR  string                    `json:"-" gorm:"perk"`
	SkillWin map[string]int            `json:"skill" gorm:"-"`
	SkillSTR string                    `json:"-" gorm:"skill"`
	SpellWin map[string]int            `json:"spell" gorm:"-"`
	SpellSTR string                    `json:"-" gorm:"spell"`
}

func (c *Champion) TableName() string {
	return "analyzed_champions"
}
func (c *Champion) MarshalBinary() ([]byte, error) {
	return json.Marshal(c)
}

func (c *Champion) UnmarshalBinary(bt []byte) error {
	return json.Unmarshal(bt, c)
}
