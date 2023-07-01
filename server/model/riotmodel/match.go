package riotmodel

import (
	"encoding/json"
	
	"gorm.io/gorm"
)

type MatchDto struct {
	gorm.Model
	Summoners []*SummonerDTO `gorm:"many2many:match_summoners"`
	
	Metadata *MetadataDto `json:"metadata" gorm:"embedded"` // 比赛元数据
	Info     *InfoDto     `json:"info" gorm:"embedded"`     // 比赛信息
}

// UnmarshalJSON assigning the matchID to each substructure
func (p *MatchDto) UnmarshalJSON(data []byte) error {
	// avoid recursion call
	var tmp struct {
		Metadata *MetadataDto `json:"metadata"` // 比赛元数据
		Info     *InfoDto     `json:"info"`     // 比赛信息
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	// assigning the matchID to each substructure
	p.Info = tmp.Info
	p.Metadata = tmp.Metadata
	matchID := p.Metadata.MetaMatchID
	for _, part := range p.Info.Participants {
		part.MetaMatchID = matchID
	}
	for _, team := range p.Info.Teams {
		team.MetaMatchID = matchID
	}
	return nil
}
