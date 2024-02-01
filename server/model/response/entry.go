package response

type EntryDTO struct {
	QueType      string `json:"queueType"`    // 排位类型
	Tier         string `json:"tier"`         // 段位
	Rank         string `json:"rank"`         // 段位
	LeaguePoints int    `json:"leaguePoints"` // 段位积分
	Wins         int    `json:"wins"`         // 胜场次数（召唤师峡谷）
	Losses       int    `json:"losses"`       // 负场次数（召唤师峡谷）
}
