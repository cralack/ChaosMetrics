package riotmodel

import (
	"encoding/json"
	"time"
)

type InfoDto struct {
	GameCreation       time.Time        `json:"gameCreation"`       // 游戏创建时间戳
	GameDuration       int64            `json:"gameDuration"`       // 游戏持续时间
	GameEndTimestamp   time.Time        `json:"gameEndTimestamp"`   // 游戏结束时间戳
	GameID             int64            `json:"gameId"`             // 游戏ID
	GameMode           string           `json:"gameMode"`           // 游戏模式
	GameName           string           `json:"gameName"`           // 游戏名称
	GameStartTimestamp time.Time        `json:"gameStartTimestamp"` // 游戏开始时间戳
	GameType           string           `json:"gameType"`           // 游戏类型
	GameVersion        string           `json:"gameVersion"`        // 游戏版本
	MapID              int              `json:"mapId"`              // 地图ID
	Participants       []ParticipantDto `json:"participants"`       // 参与者列表
	PlatformID         string           `json:"platformId"`         // 比赛所在平台
	QueueID            int              `json:"queueId"`            // 队列ID
	Teams              []TeamDto        `json:"teams"`              // 队伍列表
	TournamentCode     string           `json:"tournamentCode"`     // 生成比赛的锦标赛代码
}

type TeamDto struct {
	Bans       []BanDto      `json:"bans"`       // 禁用的英雄列表
	Objectives ObjectivesDto `json:"objectives"` // 目标
	TeamId     int           `json:"teamId"`     // 队伍ID
	Win        bool          `json:"win"`        // 是否获胜
}

type BanDto struct {
	ChampionId int `json:"championId"` // 禁用的英雄ID
	PickTurn   int `json:"pickTurn"`   // 禁用的顺序
}

type ObjectivesDto struct {
	Baron      ObjectiveDto `json:"baron"`      // 大龙
	Champion   ObjectiveDto `json:"champion"`   // 英雄
	Dragon     ObjectiveDto `json:"dragon"`     // 小龙
	Inhibitor  ObjectiveDto `json:"inhibitor"`  // 抑制塔
	RiftHerald ObjectiveDto `json:"riftHerald"` // 先驱
	Tower      ObjectiveDto `json:"tower"`      // 塔
}

type ObjectiveDto struct {
	First bool `json:"first"` // 是否是首次达成
	Kills int  `json:"kills"` // 击杀数
}

func (p *InfoDto) UnmarshalJSON(data []byte) error {
	var f map[string]interface{}
	err := json.Unmarshal(data, &f)
	if err != nil {
		return err
	}
	for k, v := range f {
		switch k {
		case "gameCreation":
			gameCreationMillis := int64(v.(float64))
			p.GameCreation = time.Unix(0, gameCreationMillis*int64(time.Millisecond)).UTC()
		case "gameDuration":
			p.GameDuration = int64(v.(float64))
		case "gameEndTimestamp":
			gameEndTimestampMillis := int64(v.(float64))
			p.GameEndTimestamp = time.Unix(0, gameEndTimestampMillis*int64(time.Millisecond)).UTC()
		case "gameId":
			p.GameID = int64(v.(float64))
		case "gameMode":
			p.GameMode = v.(string)
		case "gameName":
			p.GameName = v.(string)
		case "gameStartTimestamp":
			gameStartTimestampMillis := int64(v.(float64))
			p.GameStartTimestamp = time.Unix(0, gameStartTimestampMillis*int64(time.Millisecond)).UTC()
		case "gameType":
			p.GameType = v.(string)
		case "gameVersion":
			p.GameVersion = v.(string)
		case "mapId":
			p.MapID = int(v.(float64))
		case "participants":
			participantsJSON, err := json.Marshal(v)
			if err != nil {
				return err
			}
			err = json.Unmarshal(participantsJSON, &p.Participants)
			if err != nil {
				return err
			}
		case "platformId":
			p.PlatformID = v.(string)
		case "queueId":
			p.QueueID = int(v.(float64))
		case "teams":
			teamsJSON, err := json.Marshal(v)
			if err != nil {
				return err
			}
			err = json.Unmarshal(teamsJSON, &p.Teams)
			if err != nil {
				return err
			}
		case "tournamentCode":
			p.TournamentCode = v.(string)
		default:
			// Handle other fields as needed
		}
	}

	return nil
}
