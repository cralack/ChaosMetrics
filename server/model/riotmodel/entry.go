package riotmodel

import (
	"encoding/json"

	"gorm.io/gorm"
)

type LeagueListDTO struct {
	LeagueID string            `json:"leagueId"` // 联赛ID
	Entries  []*LeagueEntryDTO `json:"entries"`  // 联赛项列表
	Tier     string            `json:"tier"`     // 段位
	Name     string            `json:"name"`     // 联赛名称
	Queue    string            `json:"queue"`    // 游戏队列
}

type LeagueEntryDTO struct {
	gorm.Model

	Loc          string         `json:"loc" gorm:"column:loc;type:varchar(100)"`                          // 地区
	FreshBlood   bool           `json:"freshBlood" gorm:"column:fresh_blood"`                             // 是否是新晋选手
	Wins         int            `json:"wins" gorm:"column:wins;type:smallint"`                            // 胜场次数（召唤师峡谷）
	SummonerName string         `json:"summonerName" gorm:"column:summoner_name;type:varchar(100);index"` // 召唤师名称
	MiniSeries   *MiniSeriesDTO `json:"miniSeries" gorm:"embedded;embeddedPrefix:mini_"`                  // 小系列赛信息
	Inactive     bool           `json:"inactive" gorm:"column:inactive"`                                  // 是否处于非活跃状态
	Veteran      bool           `json:"veteran" gorm:"column:veteran"`                                    // 是否是资深选手
	HotStreak    bool           `json:"hotStreak" gorm:"column:hot_streak"`                               // 是否处于连胜状态
	Tier         string         `json:"tier" gorm:"column:tier;type:varchar(100)"`                        // 段位
	Rank         string         `json:"rank" gorm:"column:rank;type:varchar(100)"`                        // 段位
	LeaguePoints int            `json:"leaguePoints" gorm:"column:league_points;type:smallint"`           // 段位积分
	Losses       int            `json:"losses" gorm:"column:losses;type:smallint"`                        // 负场次数（召唤师峡谷）
	Puuid        string         `json:"puuid" gorm:"column:puuid;type:varchar(100)"`                      // 玩家的加密召唤师ID
	// riot update on jul 20,23
	LeagueID string `json:"leagueId" gorm:"column:league_id;type:varchar(100)"`   // 联赛ID
	QueType  string `json:"queueType" gorm:"column:queue_type;type:varchar(100)"` // 排位类型
}

type MiniSeriesDTO struct {
	Losses   int    `json:"losses" gorm:"column:losses"`     // 小系列赛中的失败场次
	Progress string `json:"progress" gorm:"column:progress"` // 小系列赛的进度
	Target   int    `json:"target" gorm:"column:target"`     // 小系列赛的目标胜场次
	Wins     int    `json:"wins" gorm:"column:wins"`         // 小系列赛中的胜场次
}

func (e *LeagueEntryDTO) TableName() string {
	return "entries"
}

func (e *LeagueEntryDTO) MarshalBinary() ([]byte, error)  { return json.Marshal(e) }
func (e *LeagueEntryDTO) UnmarshalBinary(bt []byte) error { return json.Unmarshal(bt, e) }
