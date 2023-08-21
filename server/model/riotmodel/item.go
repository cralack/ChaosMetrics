package riotmodel

import (
	"encoding/json"
	
	"github.com/cralack/ChaosMetrics/server/model"
)

type ItemList struct {
	Type    string              `json:"type"`
	Version string              `json:"version"`
	Basic   BasicInfo           `json:"basic"`
	Data    map[string]*ItemDTO `json:"data"`
	Groups  []*GroupInfo        `json:"groups"`
	Tree    []*TreeInfo         `json:"tree"`
}

type BasicInfo struct {
	Name string `json:"name"`
	Rune struct {
		IsRune bool   `json:"isrune"`
		Tier   int    `json:"tier"`
		Type   string `json:"type"`
	} `json:"rune"`
	Gold             *GoldInfo          `json:"gold"`
	Group            string             `json:"group"`
	Description      string             `json:"description"`
	Colloq           string             `json:"colloq"`
	Plaintext        string             `json:"plaintext"`
	Consumed         bool               `json:"consumed"`
	Stacks           int                `json:"stacks"`
	Depth            int                `json:"depth"`
	ConsumeOnFull    bool               `json:"consumeOnFull"`
	From             []string           `json:"from"`
	Into             []string           `json:"into"`
	SpecialRecipe    int                `json:"specialRecipe"`
	InStore          bool               `json:"inStore"`
	HideFromAll      bool               `json:"hideFromAll"`
	RequiredChampion string             `json:"requiredChampion"`
	RequiredAlly     string             `json:"requiredAlly"`
	Stats            map[string]float64 `json:"stats"`
	Tags             []string           `json:"tags"`
	Maps             map[string]bool    `json:"maps"`
}

type GoldInfo struct {
	Base        int  `json:"base"`
	Total       int  `json:"total"`
	Sell        int  `json:"sell"`
	Purchasable bool `json:"purchasable"`
}

type ItemDTO struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Colloq      string             `json:"colloq"`
	Plaintext   string             `json:"plaintext"`
	From        []string           `json:"from"`
	Into        []string           `json:"into"`
	Image       *model.Image       `json:"image"`
	Gold        *GoldInfo          `json:"gold"`
	Tags        []string           `json:"tags"`
	Maps        map[string]bool    `json:"maps"`
	Stats       map[string]float64 `json:"stats"`
	Depth       int                `json:"depth"`
}

func (i *ItemDTO) MarshalBinary() ([]byte, error) {
	return json.Marshal(i)
}

func (i *ItemDTO) UnmarshalBinary(bt []byte) error {
	return json.Unmarshal(bt, i)
}

type GroupInfo struct {
	ID              string `json:"id"`
	MaxGroupOwnable string `json:"MaxGroupOwnable"`
}

type TreeInfo struct {
	Header string   `json:"header"`
	Tags   []string `json:"tags"`
}
