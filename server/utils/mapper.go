package utils

import (
	"github.com/cralack/ChaosMetrics/server/model/response"
	"github.com/cralack/ChaosMetrics/server/model/riotmodel"
)

func ConvertSummonerToDTO(src *riotmodel.SummonerDTO) *response.SummonerDTO {
	if src == nil {
		return nil
	}

	return &response.SummonerDTO{
		Name:          src.Name,
		Loc:           src.Loc,
		ProfileIconID: src.ProfileIconID,
		SummonerLevel: src.SummonerLevel,
	}
}
