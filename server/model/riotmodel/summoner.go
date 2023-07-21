package riotmodel

import (
	"encoding"
	"encoding/json"
	"time"
	
	"github.com/cralack/ChaosMetrics/server/global"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SummonerDTO struct {
	gorm.Model
	Matches []*MatchDto `gorm:"many2many:match_summoners"` // 比赛列表，多对多关系 FIFO
	
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

// check implement
var _ DTO = &SummonerDTO{}

func (p *SummonerDTO) UnmarshalJSON(data []byte) error {
	var f map[string]interface{}
	
	err := json.Unmarshal(data, &f)
	if err != nil {
		return err
	}
	for k, v := range f {
		switch k {
		// gorm.Model
		case "ID":
			p.ID = uint(v.(float64))
		case "DeletedAt":
			if v != nil {
				p.DeletedAt = v.(gorm.DeletedAt)
			}
		case "CreatedAt":
			if p.CreatedAt, err = convertTime(v, time.RFC3339); err != nil {
				global.GVA_LOG.Error("parse failed", zap.Error(err))
				return err
			}
		case "UpdatedAt":
			if p.UpdatedAt, err = convertTime(v, time.RFC3339); err != nil {
				global.GVA_LOG.Error("parse failed", zap.Error(err))
				return err
			}
		
		case "accountId":
			p.AccountID = v.(string)
		case "profileIconId":
			p.ProfileIconID = int(v.(float64))
		case "revisionDate":
			// Assuming revisionDate is in milliseconds
			if revisionDateMillis, ok := v.(float64); ok {
				p.RevisionDate = time.Unix(int64(revisionDateMillis)/1000, 0).UTC()
			} else if p.RevisionDate, err = convertTime(v, time.RFC3339); err != nil {
				global.GVA_LOG.Error("parse failed", zap.Error(err))
				return err
			}
		case "name":
			p.Name = v.(string)
		case "id":
			p.MetaSummonerID = v.(string)
		case "puuid":
			p.PUUID = v.(string)
		case "summonerLevel":
			p.SummonerLevel = int(v.(float64))
			
		}
	}
	
	return nil
}

var _ encoding.BinaryMarshaler = &SummonerDTO{}

func (p *SummonerDTO) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

var _ encoding.BinaryUnmarshaler = &SummonerDTO{}

func (p *SummonerDTO) UnmarshalBinary(bt []byte) error {
	return json.Unmarshal(bt, p)
}