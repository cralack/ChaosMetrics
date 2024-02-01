package riotmodel

import (
	"encoding/json"

	"github.com/cralack/ChaosMetrics/server/model"
)

type ChampionSingleDTO struct {
	Type    string                  `json:"type"`
	Format  string                  `json:"format"`
	Version string                  `json:"version"`
	Data    map[string]*ChampionDTO `json:"data"`
}
type ChampionDTO struct {
	ID        string         `json:"id"`        // 英雄ID:Aatrox
	Key       string         `json:"key"`       // 英雄Key:266
	Name      string         `json:"name"`      // 英雄名称:暗裔剑魔
	Title     string         `json:"title"`     // 英雄称号:亚托克斯
	Image     *model.Image   `json:"image"`     // 图像
	Skins     []*Skin        `json:"skins"`     // 皮肤
	Lore      string         `json:"lore"`      // 英雄背景故事
	Blurb     string         `json:"blurb"`     // 英雄简介
	AllyTips  []string       `json:"allytips"`  // 盟友提示
	EnemyTips []string       `json:"enemytips"` // 敌人提示
	Tags      []string       `json:"tags"`      // 标签
	Partype   string         `json:"partype"`   // 资源类型
	Info      *Info          `json:"info"`      // 信息
	Stats     *ChampionStats `json:"stats"`     // 统计数据
	Spells    []*Spell       `json:"spells"`    // 技能
	Passive   struct {
		Name        string       `json:"name"`        // 被动技能名称
		Description string       `json:"description"` // 被动技能描述
		Image       *model.Image `json:"image"`       // 被动技能图像
	} `json:"passive"` // 被动技能
	Recommended []interface{} // 推荐
}

type ChampionListDTO struct {
	Type    string                      `json:"type"`
	Format  string                      `json:"format"`
	Version string                      `json:"version"`
	Data    map[string]*ChampionMiniDTO `json:"data"`
}

type ChampionMiniDTO struct {
	Version string         `json:"version"` // 版本
	ID      string         `json:"id"`      // ID:Aatrox
	Key     string         `json:"key"`     // 键:266
	Name    string         `json:"name"`    // 名称:暗裔剑魔
	Title   string         `json:"title"`   // 标题:亚托克斯
	Blurb   string         `json:"blurb"`   // 简介
	Info    *Info          `json:"info"`    // 信息
	Image   *model.Image   `json:"image"`   // 图像
	Tags    []string       `json:"tags"`    // 标签
	GTags   string         `json:"gtags"`   // db存储标签
	Partype string         `json:"partype"` // 资源类型
	Stats   *ChampionStats `json:"stats"`   // 统计数据
}

type ChampionStats struct {
	HP                   float64 `json:"hp"`                   // 生命值
	HPPerLevel           float64 `json:"hpperlevel"`           // 每级生命值增加
	MP                   float64 `json:"mp"`                   // 法力值
	MPPerLevel           float64 `json:"mpperlevel"`           // 每级法力值增加
	MoveSpeed            float64 `json:"movespeed"`            // 移动速度
	Armor                float64 `json:"armor"`                // 护甲值
	ArmorPerLevel        float64 `json:"armorperlevel"`        // 每级护甲值增加
	SpellBlock           float64 `json:"spellblock"`           // 魔法抗性值
	SpellBlockPerLevel   float64 `json:"spellblockperlevel"`   // 每级魔法抗性值增加
	AttackRange          float64 `json:"attackrange"`          // 攻击范围
	HPRegen              float64 `json:"hpregen"`              // 生命值回复
	HPRegenPerLevel      float64 `json:"hpregenperlevel"`      // 每级生命值回复增加
	MPRegen              float64 `json:"mpregen"`              // 法力值回复
	MPRegenPerLevel      float64 `json:"mpregenperlevel"`      // 每级法力值回复增加
	Crit                 float64 `json:"crit"`                 // 暴击几率
	CritPerLevel         float64 `json:"critperlevel"`         // 每级暴击几率增加
	AttackDamage         float64 `json:"attackdamage"`         // 攻击力
	AttackDamagePerLevel float64 `json:"attackdamageperlevel"` // 每级攻击力增加
	AttackSpeedPerLevel  float64 `json:"attackspeedperlevel"`  // 每级攻击速度增加
	AttackSpeed          float64 `json:"attackspeed"`          // 攻击速度
}

type Info struct {
	Attack     int `json:"attack"`     // 攻击力
	Defense    int `json:"defense"`    // 防御力
	Magic      int `json:"magic"`      // 法术强度
	Difficulty int `json:"difficulty"` // 难度等级
}

type Skin struct {
	SkinID  string `json:"id"`      // 皮肤ID
	Num     int    `json:"num"`     // 皮肤编号
	Name    string `json:"name"`    // 皮肤名称
	Chromas bool   `json:"chromas"` // 是否有多彩皮肤
}

type Spell struct {
	SpellID      string       `json:"id"`           // 技能ID:AatroxR
	Name         string       `json:"name"`         // 技能名称:大灭
	Description  string       `json:"description"`  // 技能描述
	Tooltip      string       `json:"tooltip"`      // 技能提示
	LevelTip     LevelTip     `json:"leveltip"`     // 技能等级提示
	MaxRank      int          `json:"maxrank"`      // 技能最大等级
	Cooldown     []float64    `json:"cooldown"`     // 冷却时间
	CooldownBurn string       `json:"cooldownBurn"` // 冷却时间描述
	Cost         []int        `json:"cost"`         // 花费
	CostBurn     string       `json:"costBurn"`     // 花费描述
	Range        []int        `json:"range"`        // 施法范围
	RangeBurn    string       `json:"rangeBurn"`    // 施法范围描述
	Image        *model.Image `json:"image"`        // 图像
	Resource     string       `json:"resource"`     // 资源
}

func (c *ChampionDTO) MarshalBinary() ([]byte, error)  { return json.Marshal(c) }
func (c *ChampionDTO) UnmarshalBinary(bt []byte) error { return json.Unmarshal(bt, c) }

type LevelTip struct {
	Label  []string `json:"label" gorm:"column:label"`   // 等级提示标签
	Effect []string `json:"effect" gorm:"column:effect"` // 等级提示效果
}
