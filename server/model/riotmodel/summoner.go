package riotmodel

import (
	"encoding/json"
	"time"

	"github.com/cralack/ChaosMetrics/server/internal/global"
	"github.com/cralack/ChaosMetrics/server/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SummonerDTO struct {
	gorm.Model
	Matches string `gorm:"column:matches"`

	Loc            string    `gorm:"column:loc;type:varchar(100)" json:"loc"`
	AccountID      string    `gorm:"column:account_id;type:varchar(100)" json:"accountId"`      // 加密的账号ID，最长为56个字符
	ProfileIconID  int       `gorm:"column:profile_icon_id;type:smallint" json:"profileIconId"` // 与召唤师相关联的召唤师图标ID
	RevisionDate   time.Time `gorm:"column:revision_date" json:"revisionDate"`                  // 召唤师最后修改的日期，以毫秒为单位的时间戳
	Name           string    `gorm:"column:name;index;type:varchar(100)" json:"name"`           // 召唤师名称
	MetaSummonerID string    `gorm:"column:meta_summoner_id;index;type:varchar(100)" json:"id"` // 加密的召唤师ID，最长为63个字符
	PUUID          string    `gorm:"column:puuid;type:varchar(100)" json:"puuid"`               // 加密的PUUID，长度为78个字符
	SummonerLevel  int       `gorm:"column:summoner_level;type:smallint" json:"summonerLevel"`  // 召唤师等级
	TournamentID   string    `gorm:"column:tournament_id;type:varchar(100)"`                    // 锦标赛ID
}

func (s *SummonerDTO) TableName() string {
	return "summoners"
}

func (s *SummonerDTO) UnmarshalJSON(data []byte) error {
	var f map[string]interface{}

	err := json.Unmarshal(data, &f)
	if err != nil {
		return err
	}
	for k, v := range f {
		switch k {
		// gorm.Model
		case "ID":
			s.ID = uint(v.(float64))
		case "DeletedAt":
			if v != nil {
				s.DeletedAt = v.(gorm.DeletedAt)
			}
		case "CreatedAt":
			if s.CreatedAt, err = model.ConvertTime(v, time.RFC3339); err != nil {
				global.GvaLog.Error("parse failed", zap.Error(err))
				return err
			}
		case "UpdatedAt":
			if s.UpdatedAt, err = model.ConvertTime(v, time.RFC3339); err != nil {
				global.GvaLog.Error("parse failed", zap.Error(err))
				return err
			}

		case "accountId":
			s.AccountID = v.(string)
		case "profileIconId":
			s.ProfileIconID = int(v.(float64))
		case "revisionDate":
			// Assuming revisionDate is in milliseconds
			if revisionDateMillis, ok := v.(float64); ok {
				s.RevisionDate = time.Unix(int64(revisionDateMillis)/1000, 0).UTC()
			} else if s.RevisionDate, err = model.ConvertTime(v, time.RFC3339); err != nil {
				global.GvaLog.Error("parse failed", zap.Error(err))
				return err
			}
		case "name":
			s.Name = v.(string)
		case "id":
			s.MetaSummonerID = v.(string)
		case "puuid":
			s.PUUID = v.(string)
		case "summonerLevel":
			s.SummonerLevel = int(v.(float64))

		}
	}

	return nil
}

func (s *SummonerDTO) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}
func (s *SummonerDTO) UnmarshalBinary(bt []byte) error {
	return json.Unmarshal(bt, s)
}
