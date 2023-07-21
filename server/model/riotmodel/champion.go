package riotmodel

type ChampionSingleDTO struct {
	Type    string                 `json:"type"`
	Format  string                 `json:"format"`
	Version string                 `json:"version"`
	Data    map[string]ChampionDTO `json:"data"`
}

type ChampionDTO struct {
	version   string
	ID        string         `json:"id"`
	Key       uint           `json:"key"`
	Name      string         `json:"name"`
	Title     string         `json:"title"`
	Image     *Image         `json:"image"`
	Skins     []*Skin        `json:"skins"`
	Lore      string         `json:"lore"`
	Blurb     string         `json:"blurb"`
	AllyTips  []string       `json:"allytips"`
	EnemyTips []string       `json:"enemytips"`
	Tags      []string       `json:"tags"`
	Partype   string         `json:"partype"`
	Info      *Info          `json:"info"`
	Stats     *ChampionStats `json:"stats"`
	Spells    []Spell        `json:"spells"`
	Passive   struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Image       *Image `json:"image"`
	} `json:"passive"`
	Recommended []interface{} `json:"recommended"`
}

type ChampionListDTO struct {
	Type    string                      `json:"type"`
	Format  string                      `json:"format"`
	Version string                      `json:"version"`
	Data    map[string]*ChampionMiniDTO `json:"data"`
}

type ChampionMiniDTO struct {
	Version string         `json:"version"`
	ID      string         `json:"id"`
	Key     string         `json:"key"`
	Name    string         `json:"name"`
	Title   string         `json:"title"`
	Blurb   string         `json:"blurb"`
	Info    *Info          `json:"info"`
	Image   *Image         `json:"image"`
	Tags    []string       `json:"tags"`
	Partype string         `json:"partype"`
	Stats   *ChampionStats `json:"stats"`
}

type ChampionStats struct {
	HP                   float64 `json:"hp"`
	HPPerLevel           float64 `json:"hpperlevel"`
	MP                   float64 `json:"mp"`
	MPPerLevel           float64 `json:"mpperlevel"`
	MoveSpeed            float64 `json:"movespeed"`
	Armor                float64 `json:"armor"`
	ArmorPerLevel        float64 `json:"armorperlevel"`
	SpellBlock           float64 `json:"spellblock"`
	SpellBlockPerLevel   float64 `json:"spellblockperlevel"`
	AttackRange          int     `json:"attackrange"`
	HPRegen              float64 `json:"hpregen"`
	HPRegenPerLevel      float64 `json:"hpregenperlevel"`
	MPRegen              float64 `json:"mpregen"`
	MPRegenPerLevel      float64 `json:"mpregenperlevel"`
	Crit                 int     `json:"crit"`
	CritPerLevel         int     `json:"critperlevel"`
	AttackDamage         float64 `json:"attackdamage"`
	AttackDamagePerLevel float64 `json:"attackdamageperlevel"`
	AttackSpeedPerLevel  float64 `json:"attackspeedperlevel"`
	AttackSpeed          float64 `json:"attackspeed"`
}
type Info struct {
	Attack     int `json:"attack"`
	Defense    int `json:"defense"`
	Magic      int `json:"magic"`
	Difficulty int `json:"difficulty"`
}

type Skin struct {
	ID      string `json:"id"`
	Num     int    `json:"num"`
	Name    string `json:"name"`
	Chromas bool   `json:"chromas"`
}

type Spell struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Tooltip      string    `json:"tooltip"`
	LevelTip     LevelTip  `json:"leveltip"`
	MaxRank      int       `json:"maxrank"`
	Cooldown     []float64 `json:"cooldown"`
	CooldownBurn string    `json:"cooldownBurn"`
	Cost         []int     `json:"cost"`
	CostBurn     string    `json:"costBurn"`
	Range        []int     `json:"range"`
	RangeBurn    string    `json:"rangeBurn"`
	Image        *Image    `json:"image"`
	Resource     string    `json:"resource"`
}

type LevelTip struct {
	Label  []string `json:"label"`
	Effect []string `json:"effect"`
}
