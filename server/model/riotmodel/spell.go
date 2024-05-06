package riotmodel

import (
	"encoding/json"

	"github.com/cralack/ChaosMetrics/server/model"
)

type SpellList struct {
	Type    string                    `json:"type"`
	Version string                    `json:"version"`
	Data    map[string]*SummonerSpell `json:"data"`
}

type SummonerSpell struct {
	ID            string       `json:"id"`
	Name          string       `json:"name"`
	Description   string       `json:"description"`
	Tooltip       string       `json:"tooltip"`
	MaxRank       int          `json:"maxrank"`
	Cooldown      []float32    `json:"cooldown"`
	CooldownBurn  string       `json:"cooldownBurn"`
	Cost          []int        `json:"cost"`
	CostBurn      string       `json:"costBurn"`
	DataValues    struct{}     `json:"datavalues"`
	Effect        [][]float32  `json:"effect"`
	EffectBurn    []string     `json:"effectBurn"`
	Vars          []struct{}   `json:"vars"`
	Key           string       `json:"key"`
	SummonerLevel int          `json:"summonerLevel"`
	Modes         []string     `json:"modes"`
	CostType      string       `json:"costType"`
	MaxAmmo       string       `json:"maxammo"`
	Range         []int        `json:"range"`
	RangeBurn     string       `json:"rangeBurn"`
	Resource      string       `json:"resource"`
	Image         *model.Image `json:"image"`
}

func (s *SpellList) MarshalBinary() ([]byte, error)  { return json.Marshal(s) }
func (s *SpellList) UnmarshalBinary(bt []byte) error { return json.Unmarshal(bt, s) }
