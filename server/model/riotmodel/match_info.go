package riotmodel

import (
	"encoding/json"
	"time"
)

type MetadataDto struct {
	DataVersion  string   `json:"dataVersion" gorm:"column:data_version" ` // 比赛数据版本
	MetaMatchID  string   `json:"matchId" gorm:"column:meta_match_id"`     // 比赛ID
	Participants []string `json:"participants" gorm:"-"`                   // 参与者 PUUID 列表
}

type InfoDto struct {
	GameCreation       time.Time         `json:"gameCreation" gorm:"column:game_creation"`              // 游戏创建时间戳
	GameDuration       int               `json:"gameDuration" gorm:"column:game_duration"`              // 游戏持续时间
	GameEndTimestamp   time.Time         `json:"gameEndTimestamp" gorm:"column:game_end_timestamp"`     // 游戏结束时间戳
	GameID             int               `json:"gameId" gorm:"column:game_id"`                          // 游戏ID
	GameMode           string            `json:"gameMode" gorm:"column:game_mode"`                      // 游戏模式
	GameName           string            `json:"gameName" gorm:"column:game_name"`                      // 游戏名称
	GameStartTimestamp time.Time         `json:"gameStartTimestamp" gorm:"column:game_start_timestamp"` // 游戏开始时间戳
	GameType           string            `json:"gameType" gorm:"column:game_type"`                      // 游戏类型
	GameVersion        string            `json:"gameVersion" gorm:"column:game_version"`                // 游戏版本
	MapID              int               `json:"mapId" gorm:"column:map_id"`                            // 地图ID
	Participants       []*ParticipantDto `json:"participants" gorm:"foreignKey:match_id"`               // 参与者列表
	PlatformID         string            `json:"platformId" gorm:"column:platform_id"`                  // 比赛所在平台
	QueueID            int               `json:"queueId" gorm:"column:queue_id"`                        // 队列ID
	Teams              []*TeamDto        `json:"teams" gorm:"-"`                                        // 队伍列表
	TournamentCode     string            `json:"tournamentCode" gorm:"column:tournament_code"`          // 生成比赛的锦标赛代码
}

// TeamDto todo:add gorm tag
type TeamDto struct {
	Bans       []*BanDto      `json:"bans"`       // 禁用的英雄列表
	Objectives *ObjectivesDto `json:"objectives"` // 目标
	TeamId     int            `json:"teamId"`     // 队伍ID
	Win        bool           `json:"win"`        // 是否获胜
}

type BanDto struct {
	ChampionId int `json:"championId"` // 禁用的英雄ID
	PickTurn   int `json:"pickTurn"`   // 禁用的顺序
}

type ObjectivesDto struct {
	Baron      *ObjectiveDto `json:"baron"`      // 大龙
	Champion   *ObjectiveDto `json:"champion"`   // 英雄
	Dragon     *ObjectiveDto `json:"dragon"`     // 小龙
	Inhibitor  *ObjectiveDto `json:"inhibitor"`  // 抑制塔
	RiftHerald *ObjectiveDto `json:"riftHerald"` // 先驱
	Tower      *ObjectiveDto `json:"tower"`      // 塔
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
			p.GameDuration = int(v.(float64))
		case "gameEndTimestamp":
			gameEndTimestampMillis := int64(v.(float64))
			p.GameEndTimestamp = time.Unix(0, gameEndTimestampMillis*int64(time.Millisecond)).UTC()
		case "gameId":
			p.GameID = int(v.(float64))
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
			for _, par := range p.Participants {
				par.PerksORM = &PerksORM{}
				par.PerksMeta.parsePerksMeta(par.PerksORM)
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
