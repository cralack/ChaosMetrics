package riotmodel

import (
	"encoding/json"
)

type Rune struct {
	ID        int    `json:"id"`
	Key       string `json:"key"`
	Icon      string `json:"icon"`
	Name      string `json:"name"`
	ShortDesc string `json:"shortDesc"`
	LongDesc  string `json:"longDesc"`
}

type Slot struct {
	Runes []*Rune `json:"runes"`
}

type Perk struct {
	ID    int     `json:"id"`
	Key   string  `json:"key"`
	Icon  string  `json:"icon"`
	Name  string  `json:"name"`
	Slots []*Slot `json:"slots"`
}

func (p *Perk) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

func (p *Perk) UnmarshalBinary(bt []byte) error {
	return json.Unmarshal(bt, p)
}

type PerkDetail struct {
	ID                                 int            `json:"id"`
	Name                               string         `json:"name"`
	MajorChangePatchVersion            string         `json:"majorChangePatchVersion"`
	Tooltip                            string         `json:"tooltip"`
	ShortDesc                          string         `json:"shortDesc"`
	LongDesc                           string         `json:"longDesc"`
	RecommendationDescriptor           string         `json:"recommendationDescriptor"`
	IconPath                           string         `json:"iconPath"`
	EndOfGameStatDescs                 []string       `json:"endOfGameStatDescs"`
	RecommendationDescriptorAttributes map[string]int `json:"recommendationDescriptorAttributes"`
}

func (p *PerkDetail) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

func (p *PerkDetail) UnmarshalBinary(bt []byte) error {
	return json.Unmarshal(bt, p)
}
