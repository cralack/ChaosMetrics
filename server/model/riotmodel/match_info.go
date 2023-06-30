package riotmodel

import (
	"encoding/json"
	"time"
	
	"gorm.io/gorm"
)

type MetadataDto struct {
	DataVersion  string   `json:"dataVersion" gorm:"column:data_version;type:varchar(100)"` // 比赛数据版本
	MetaMatchID  string   `json:"matchId" gorm:"column:meta_match_id;type:varchar(100)"`    // 比赛ID
	Participants []string `json:"participants" gorm:"-"`                                    // 参与者 PUUID 列表
}

type InfoDto struct {
	GameCreation       time.Time         `json:"gameCreation" gorm:"column:game_creation"`                       // 游戏创建时间戳
	GameDuration       int               `json:"gameDuration" gorm:"column:game_duration;type:smallint"`         // 游戏持续时间
	GameEndTimestamp   time.Time         `json:"gameEndTimestamp" gorm:"column:game_end_timestamp"`              // 游戏结束时间戳
	GameID             int               `json:"gameId" gorm:"column:game_id;type:int"`                          // 游戏ID
	GameMode           string            `json:"gameMode" gorm:"column:game_mode;type:varchar(100)"`             // 游戏模式
	GameName           string            `json:"gameName" gorm:"column:game_name;type:varchar(100)"`             // 游戏名称
	GameStartTimestamp time.Time         `json:"gameStartTimestamp" gorm:"column:game_start_timestamp"`          // 游戏开始时间戳
	GameType           string            `json:"gameType" gorm:"column:game_type;type:varchar(100)"`             // 游戏类型
	GameVersion        string            `json:"gameVersion" gorm:"column:game_version;type:varchar(100)"`       // 游戏版本
	MapID              int               `json:"mapId" gorm:"column:map_id;type:smallint"`                       // 地图ID
	Participants       []*ParticipantDto `json:"participants" gorm:"foreignKey:match_id"`                        // 参与者列表
	PlatformID         string            `json:"platformId" gorm:"column:platform_id;type:varchar(100)"`         // 比赛所在平台
	QueueID            int               `json:"queueId" gorm:"column:queue_id;type:smallint"`                   // 队列ID
	Teams              []*TeamDto        `json:"teams"  gorm:"foreignKey:match_id"`                              // 队伍列表
	TournamentCode     string            `json:"tournamentCode" gorm:"column:tournament_code;type:varchar(100)"` // 生成比赛的锦标赛代码
}

// TeamDto todo:add gorm tag
type TeamDto struct {
	gorm.Model
	MatchID     string `gorm:"column:match_id;type:varchar(100)"`            // 匹配matchDTO_id
	MetaMatchID string `gorm:"column:meta_match_id;index;type:varchar(100)"` // 比赛ID
	
	Bans       []*BanDto      `json:"bans" gorm:"-"`                       // 禁用的英雄列表
	Objectives *ObjectivesDto `json:"objectives" gorm:"embedded"`          // 目标
	TeamId     int            `json:"teamId" gorm:"team_id;type:smallint"` // 队伍ID
	Win        bool           `json:"win" gorm:"win"`                      // 是否获胜
}

type BanDto struct {
	ChampionId int `json:"championId;type:smallint"` // 禁用的英雄ID
	PickTurn   int `json:"pickTurn;type:smallint"`   // 禁用的顺序
}

type ObjectivesDto struct {
	Baron      *ObjectiveDto `json:"baron" gorm:"embedded;embeddedPrefix:baron_ "`            // 大龙
	Champion   *ObjectiveDto `json:"champion" gorm:"embedded;embeddedPrefix:champion_ "`      // 英雄
	Dragon     *ObjectiveDto `json:"dragon" gorm:"embedded;embeddedPrefix:dragon_ "`          // 小龙
	Inhibitor  *ObjectiveDto `json:"inhibitor" gorm:"embedded;embeddedPrefix:inhibitor_ "`    // 抑制塔
	RiftHerald *ObjectiveDto `json:"riftHerald" gorm:"embedded;embeddedPrefix:rift_herald_ "` // 先锋
	Tower      *ObjectiveDto `json:"tower" gorm:"embedded;embeddedPrefix:tower_ "`            // 塔
}

type ObjectiveDto struct {
	First bool `json:"first" gorm:"first"`              // 是否是首次达成
	Kills int  `json:"kills" gorm:"kill;type:smallint"` // 击杀数
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
