package response

import (
	"encoding/json"
)

type MatchDTO struct {
	MatchID      string         `json:"matchID"`
	QueueID      int            `json:"queueId"` // 队列ID
	GameCreation int64          `json:"gameCreation"`
	GameDuration int            `json:"gameDuration"`
	Participants []*Participant `json:"participants"`
}

type Participant struct {
	SummonerName string  `json:"summonerName"` // 召唤师名称
	Tagline      string  `json:"tagline"`
	ChampionName string  `json:"championName"` // 英雄名称
	Kills        int     `json:"kills"`        // 击杀数
	Deaths       int     `json:"deaths"`       // 死亡数
	Assists      int     `json:"assists"`      // 助攻数
	Item0        int     `json:"item0"`        // 物品0
	Item1        int     `json:"item1"`        // 物品1
	Item2        int     `json:"item2"`        // 物品2
	Item3        int     `json:"item3"`        // 物品3
	Item4        int     `json:"item4"`        // 物品4
	Item5        int     `json:"item5"`        // 物品5
	Item6        int     `json:"item6"`        // 物品6 (饰品)
	KDA          float32 `json:"kda"`          // KDA
	KP           float32 `json:"kp"`           // 击杀参与率
	PentaKills   int     `json:"pentaKills"`   // 五杀数
	QuadraKills  int     `json:"quadraKills"`  // 四杀数
	TripleKills  int     `json:"tripleKills"`  // 三杀数
	Summoner1Id  int     `json:"summoner1Id"`  // 召唤师技能1ID
	Summoner2Id  int     `json:"summoner2Id"`  // 召唤师技能2ID
	TeamId       int     `json:"teamId"`       // 队伍ID
	DamageDealt  int     `json:"damageDealt"`  // 造成伤害
	DamageToken  int     `json:"damageToken"`  // 承受伤害
	SkillBuild   string  `json:"skillBuild"`   // 技能构筑
	ItemBuild    string  `json:"itemBuild"`    // 出装构筑
	PerkBuild    string  `json:"perkBuild"`    // 符文构筑
}

func (m *MatchDTO) MarshalBinary() ([]byte, error) { return json.Marshal(m) }

func (m *MatchDTO) UnmarshalBinary(bt []byte) error { return json.Unmarshal(bt, m) }
