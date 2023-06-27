package riotmodel

type MatchDto struct {
	Metadata MetadataDto `json:"metadata"` // 比赛元数据
	Info     InfoDto     `json:"info"`     // 比赛信息
}

type MetadataDto struct {
	DataVersion  string   `json:"dataVersion"`  // 比赛数据版本
	MatchID      string   `json:"matchId"`      // 比赛ID
	Participants []string `json:"participants"` // 参与者 PUUID 列表
}
