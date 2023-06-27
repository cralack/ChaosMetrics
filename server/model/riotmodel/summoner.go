package riotmodel

import (
	"encoding/json"
	"time"
)

type SummonerDTO struct {
	AccountID     string    `json:"accountId"`     // 加密的账号ID，最长为56个字符
	ProfileIconID int       `json:"profileIconId"` // 与召唤师相关联的召唤师图标ID
	RevisionDate  time.Time `json:"revisionDate"`  // 召唤师最后修改的日期，以毫秒为单位的时间戳
	Name          string    `json:"name"`          // 召唤师名称
	ID            string    `json:"id"`            // 加密的召唤师ID，最长为63个字符
	PUUID         string    `json:"puuid"`         // 加密的PUUID，长度为78个字符
	SummonerLevel int       `json:"summonerLevel"` // 召唤师等级
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
		case "accountId":
			p.AccountID = v.(string)
		case "profileIconId":
			p.ProfileIconID = int(v.(float64))
		case "revisionDate":
			// Assuming revisionDate is in milliseconds
			revisionDateMillis := int64(v.(float64))
			p.RevisionDate = time.Unix(revisionDateMillis/1000, 0).UTC()
		case "name":
			p.Name = v.(string)
		case "id":
			p.ID = v.(string)
		case "puuid":
			p.PUUID = v.(string)
		case "summonerLevel":
			p.SummonerLevel = int(v.(float64))
		}
	}

	return nil
}
