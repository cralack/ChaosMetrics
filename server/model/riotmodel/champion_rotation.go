package riotmodel

import (
	"encoding/json"
)

type ChampionRotationDTO struct {
	FreeChampionIds              []int `json:"freeChampionIds"`
	FreeChampionIdsForNewPlayers []int `json:"freeChampionIdsForNewPlayers"`
	MaxNewPlayerLevel            int   `json:"maxNewPlayerLevel"`
}

func (p *ChampionRotationDTO) UnmarshalJSON(data []byte) error {
	var aux struct {
		FreeChampionIds              []int `json:"freeChampionIds"`
		FreeChampionIdsForNewPlayers []int `json:"freeChampionIdsForNewPlayers"`
		MaxNewPlayerLevel            int   `json:"maxNewPlayerLevel"`
	}
	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}

	p.FreeChampionIds = aux.FreeChampionIds
	p.FreeChampionIdsForNewPlayers = aux.FreeChampionIdsForNewPlayers
	p.MaxNewPlayerLevel = aux.MaxNewPlayerLevel

	return nil
}
