package riotmodel

import (
	"encoding/json"
)

type MatchDTO struct {
	Metadata *MetadataDto `json:"metadata" gorm:"embedded"` // 比赛元数据
	Info     *InfoDto     `json:"info" gorm:"embedded"`     // 比赛信息
}

// UnmarshalJSON assigning the matchID to each substructure
func (p *MatchDTO) UnmarshalJSON(data []byte) error {
	// avoid recursion call
	var tmp struct {
		Metadata *MetadataDto `json:"metadata"` // 比赛元数据
		Info     *InfoDto     `json:"info"`     // 比赛信息
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	p.Info = tmp.Info
	p.Metadata = tmp.Metadata
	return nil
}
