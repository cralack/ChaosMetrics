package riotmodel

import (
	"encoding/json"
)

type MatchDto struct {
	Metadata *MetadataDto `json:"metadata"` // 比赛元数据
	Info     *InfoDto     `json:"info"`     // 比赛信息
}

type MetadataDto struct {
	DataVersion  string   `gorm:"column:data_version" json:"dataVersion"`    // 比赛数据版本
	MatchID      string   `gorm:"primaryKey;column:match_id" json:"matchId"` // 比赛ID
	Participants []string `gorm:"-" json:"participants"`                     // 参与者 PUUID 列表
}

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
	matchID := p.Metadata.MatchID
	p.Info.MatchID = matchID
	parts := p.Info.Participants
	for _, part := range parts {
		part.MatchID = matchID
	}
	return nil
}
