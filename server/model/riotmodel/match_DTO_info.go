package riotmodel

import (
	"encoding/json"
	"time"
	
	"github.com/cralack/ChaosMetrics/server/global"
	"go.uber.org/zap"
)

type MetadataDto struct {
	DataVersion  string   `json:"dataVersion" gorm:"column:data_version;type:varchar(100)"`    // 比赛数据版本
	MetaMatchID  string   `json:"matchId" gorm:"column:meta_match_id;index;type:varchar(100)"` // 比赛ID
	Participants []string `json:"participants" gorm:"-"`                                       // 参与者 PUUID 列表
}

type InfoDto struct {
	Loc                string         `json:"loc" gorm:"column:loc;type:varchar(20)"`                         // 游戏所在服务器
	GameCreation       time.Time      `json:"gameCreation" gorm:"column:game_creation"`                       // 游戏创建时间戳
	GameDuration       int            `json:"gameDuration" gorm:"column:game_duration;type:smallint"`         // 游戏持续时间
	GameEndTimestamp   time.Time      `json:"gameEndTimestamp" gorm:"column:game_end_timestamp"`              // 游戏结束时间戳
	GameID             int            `json:"gameId" gorm:"column:game_id;type:int"`                          // 游戏ID
	GameMode           string         `json:"gameMode" gorm:"column:game_mode;type:varchar(100)"`             // 游戏模式
	GameName           string         `json:"gameName" gorm:"column:game_name;type:varchar(100)"`             // 游戏名称
	GameStartTimestamp time.Time      `json:"gameStartTimestamp" gorm:"column:game_start_timestamp"`          // 游戏开始时间戳
	GameType           string         `json:"gameType" gorm:"column:game_type;type:varchar(100)"`             // 游戏类型
	GameVersion        string         `json:"gameVersion" gorm:"column:game_version;type:varchar(100)"`       // 游戏版本
	MapID              int            `json:"mapId" gorm:"column:map_id;type:smallint"`                       // 地图ID
	Participants       []*Participant `json:"participants" gorm:"foreignKey:match_id"`                        // 参与者列表
	PlatformID         string         `json:"platformId" gorm:"column:platform_id;type:varchar(100)"`         // 比赛所在平台
	QueueID            int            `json:"queueId" gorm:"column:queue_id;type:smallint"`                   // 队列ID
	Teams              []*Team        `json:"teams"  gorm:"foreignKey:match_id"`                              // 队伍列表
	TournamentCode     string         `json:"tournamentCode" gorm:"column:tournament_code;type:varchar(100)"` // 生成比赛的锦标赛代码
}

type Team struct {
	MatchID     uint   `gorm:"column:match_id"`                              // 匹配matchDTO_id
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
	layout := time.RFC3339
	var f map[string]interface{}
	err := json.Unmarshal(data, &f)
	if err != nil {
		return err
	}
	for k, v := range f {
		switch k {
		case "gameCreation":
			if gameCreationMillis, ok := v.(float64); ok {
				p.GameCreation = time.Unix(0, int64(gameCreationMillis)*int64(time.Millisecond)).UTC()
			}
			if gameCreationMillis, ok := v.(string); ok {
				if p.GameCreation, err = time.Parse(layout, gameCreationMillis); err != nil {
					global.GVA_LOG.Error("parse failed", zap.Error(err))
					return err
				}
			}
		case "gameDuration":
			p.GameDuration = int(v.(float64))
		case "gameEndTimestamp":
			if gameEndTimestampMillis, ok := v.(float64); ok {
				p.GameEndTimestamp = time.Unix(0, int64(gameEndTimestampMillis)*int64(time.Millisecond)).UTC()
			}
			if gameEndTimestampMillis, ok := v.(string); ok {
				if p.GameEndTimestamp, err = time.Parse(layout, gameEndTimestampMillis); err != nil {
					global.GVA_LOG.Error("parse failed", zap.Error(err))
					return err
				}
			}
		case "gameId":
			p.GameID = int(v.(float64))
		case "gameMode":
			p.GameMode = v.(string)
		case "gameName":
			p.GameName = v.(string)
		case "gameStartTimestamp":
			if gameStartTimestampMillis, ok := v.(float64); ok {
				p.GameStartTimestamp = time.Unix(0, int64(gameStartTimestampMillis)*int64(time.Millisecond)).UTC()
			}
			if gameStartTimestampMillis, ok := v.(string); ok {
				if p.GameStartTimestamp, err = time.Parse(layout, gameStartTimestampMillis); err != nil {
					global.GVA_LOG.Error("parse failed", zap.Error(err))
					return err
				}
			}
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
				par.PerksSTR = par.PerksMeta.parsePerksMeta()
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
