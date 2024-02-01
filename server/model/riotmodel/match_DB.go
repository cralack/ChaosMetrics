package riotmodel

import (
	"encoding/json"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type MatchDB struct {
	gorm.Model

	Analyzed       bool             `gorm:"column:analyzed;not null;default:false"`                      // 分析标记
	MetaMatchID    string           `json:"matchId" gorm:"column:meta_match_id;index;type:varchar(100)"` // 比赛ID
	Loc            string           `json:"loc" gorm:"column:loc;type:varchar(20)"`                      // 游戏所在服务器
	Bans           string           `json:"-" gorm:"column:bans"`
	GameCreation   time.Time        `json:"gameCreation" gorm:"column:game_creation"`                       // 游戏创建时间戳
	GameDuration   int              `json:"gameDuration" gorm:"column:game_duration;type:smallint"`         // 游戏持续时间
	GameMode       string           `json:"gameMode" gorm:"column:game_mode;type:varchar(100)"`             // 游戏模式
	GameVersion    string           `json:"gameVersion" gorm:"column:game_version;type:varchar(100)"`       // 游戏版本
	MapID          int              `json:"mapId" gorm:"column:map_id;type:smallint"`                       // 地图ID
	QueueID        int              `json:"queueId" gorm:"column:queue_id;type:smallint"`                   // 队列ID (420:rank,430:match,450:ARAM,1700:CHERRY)
	Participants   []*ParticipantDB `json:"participants" gorm:"foreignKey:match_id"`                        // 参与者列表
	TournamentCode string           `json:"tournamentCode" gorm:"column:tournament_code;type:varchar(100)"` // 生成比赛的锦标赛代码
}

func (m *MatchDB) TableName() string {
	return "matches"
}

type ParticipantDB struct {
	gorm.Model

	MatchID                   uint    `json:"-" gorm:"column:match_id"`                                              // 匹配matchDTO_id
	MetaMatchID               string  `json:"match_id" gorm:"column:meta_match_id;index;type:varchar(100)"`          // 比赛ID
	Kills                     int     `json:"kills" gorm:"column:kills;type:smallint"`                               // 击杀数
	Deaths                    int     `json:"deaths" gorm:"column:deaths;type:smallint"`                             // 死亡数
	Assists                   int     `json:"assists" gorm:"column:assists;type:smallint"`                           // 助攻数
	Item0                     int     `json:"item0" gorm:"column:item0;type:int"`                                    // 物品0
	Item1                     int     `json:"item1" gorm:"column:item1;type:int"`                                    // 物品1
	Item2                     int     `json:"item2" gorm:"column:item2;type:int"`                                    // 物品2
	Item3                     int     `json:"item3" gorm:"column:item3;type:int"`                                    // 物品3
	Item4                     int     `json:"item4" gorm:"column:item4;type:int"`                                    // 物品4
	Item5                     int     `json:"item5" gorm:"column:item5;type:int"`                                    // 物品5
	Item6                     int     `json:"item6" gorm:"column:item6;type:int"`                                    // 物品6 (饰品)
	PentaKills                int     `json:"pentaKills" gorm:"column:penta_kills;type:smallint"`                    // 五杀数
	QuadraKills               int     `json:"quadraKills" gorm:"column:quadra_kills;type:smallint"`                  // 四杀数
	TripleKills               int     `json:"tripleKills" gorm:"column:triple_kills;type:smallint"`                  // 三杀数
	ProfileIcon               int     `json:"profileIcon" gorm:"column:profile_icon;type:int"`                       // 头像图标ID
	Role                      string  `json:"role" gorm:"column:role;type:varchar(100)"`                             // 角色
	ChampionName              string  `json:"championName" gorm:"column:champion_name;type:varchar(50)"`             // 英雄名称
	ChampionID                string  `json:"championId" gorm:"column:champion_id"`                                  // 英雄ID
	ChampLevel                int     `json:"champLevel" gorm:"column:champ_level;type:smallint"`                    // 英雄等级
	GameEndedInEarlySurrender bool    `json:"gameEndedInEarlySurrender" gorm:"column:game_ended_in_early_surrender"` // 是否提前投降结束游戏
	Puuid                     string  `json:"puuid" gorm:"column:puuid;type:varchar(100)"`                           // 参与者UUID
	SummonerName              string  `json:"summonerName" gorm:"column:summoner_name;type:varchar(100)"`            // 召唤师名称
	Summoner1Id               int     `json:"summoner1Id" gorm:"column:summoner1_id;type:smallint"`                  // 召唤师技能1ID
	Summoner2Id               int     `json:"summoner2Id" gorm:"column:summoner2_id;type:smallint"`                  // 召唤师技能2ID
	JudgeScore                float32 `json:"judgeScore" gorm:"column:judge_score"`                                  // 评分
	KDA                       float32 `json:"kda" gorm:"column:kda"`                                                 // KDA
	KP                        float32 `json:"kp" gorm:"column:kp"`                                                   // 击杀参与率
	TeamId                    int     `json:"teamId" gorm:"column:team_id;type:smallint"`                            // 队伍ID
	DamageDealt               int     `json:"damageDealt" gorm:"column:damage_dealt"`                                // 造成伤害
	DamageToken               int     `json:"damageToken" gorm:"column:damage_token"`                                // 承受伤害
	VisionScore               int     `json:"visionScore" gorm:"column:vision_score;type:smallint"`                  // 视野得分
	TimeCCingOthers           int     `json:"timeCCingOthers" gorm:"column:time_ccing_others;type:smallint"`         // 控制敌方英雄时间
	TotalTimeSpentDead        int     `json:"totalTimeSpentDead" gorm:"column:total_time_spent_dead;type:smallint"`  // 总死亡时间
	TotalMinionsKilled        int     `json:"totalMinionsKilled" gorm:"column:total_minions_killed;type:smallint"`   // 补刀总数
	Build                     *Build  `gorm:"embedded;embeddedPrefix:build_"`
	Win                       bool    `json:"win" gorm:"column:win"` // 是否获胜
}

func (p *ParticipantDB) TableName() string {
	return "match_participants"
}

type Build struct {
	Perk       string       `json:"perk" gorm:"column:perk"`   // 天赋
	SkillOrder []int        `json:"skillOrder" gorm:"-"`       // 加点顺序
	Skill      string       `json:"skill" gorm:"column:skill"` // 技能
	ItemOrder  []*ItemOrder `json:"itemOrder" gorm:"-"`        // 出装顺序
	Item       string       `json:"item" gorm:"column:item"`   // 物品
}

type ItemOrder struct {
	ItemID    int           `json:"itemID" gorm:"column:item_id"`       // 物品ID
	TimeStamp time.Duration `json:"timeStamp" gorm:"column:time_stamp"` // 时间戳
	Sold      bool          `json:"sold" gorm:"column:sold"`            // 购买/售出
}

func (m *MatchDB) ParseClassicAndARAM(match *MatchDTO, matchTL *MatchTimelineDTO) error {
	// parse to result
	partySize := len(match.Info.Participants) // raw json start from 1
	partics := make([]*ParticipantDB, partySize)
	if match.Info.GameMode == "CLASSIC" {
		bans := make([]int, 0, partySize)
		for _, t := range match.Info.Teams {
			for _, b := range t.Bans {
				bans = append(bans, b.ChampionId)
			}
		}
		if buff, err := json.Marshal(bans); err != nil {
			return err
		} else {
			m.Bans = string(buff)
		}
	}
	for _, p := range match.Info.Participants {
		if p.Challenges == nil {
			return nil
		}
		partics[p.ParticipantId-1] = &ParticipantDB{
			MetaMatchID:               match.Metadata.MatchID,
			Assists:                   p.Assists,
			Deaths:                    p.Deaths,
			Kills:                     p.Kills,
			Item0:                     p.Item0,
			Item1:                     p.Item1,
			Item2:                     p.Item2,
			Item3:                     p.Item3,
			Item4:                     p.Item4,
			Item5:                     p.Item5,
			Item6:                     p.Item6,
			PentaKills:                p.PentaKills,
			QuadraKills:               p.QuadraKills,
			TripleKills:               p.TripleKills,
			ProfileIcon:               p.ProfileIcon,
			Role:                      p.Role,
			ChampionName:              p.ChampionName,
			ChampionID:                strconv.Itoa(p.ChampionId),
			ChampLevel:                p.ChampLevel,
			GameEndedInEarlySurrender: p.GameEndedInEarlySurrender,
			Puuid:                     p.Puuid,
			SummonerName:              p.SummonerName,
			Summoner1Id:               p.Summoner1Id,
			Summoner2Id:               p.Summoner2Id,
			KDA:                       p.Challenges.KDA,
			KP:                        p.Challenges.KillParticipation,
			TeamId:                    p.TeamId,
			DamageDealt:               p.TotalDamageDealtToChampions,
			DamageToken:               p.TotalDamageTaken,
			VisionScore:               p.VisionScore,
			TimeCCingOthers:           p.TimeCCingOthers,
			TotalTimeSpentDead:        p.TotalTimeSpentDead,
			TotalMinionsKilled:        p.TotalMinionsKilled,
			Build: &Build{
				Perk: p.PerksSTR,
			},
			Win: p.Win,
		}
	}
	// timeline parse val
	skillSlot := make([][]int, partySize)
	itemOrder := make([][]*ItemOrder, partySize)
	for i := range skillSlot {
		skillSlot[i] = make([]int, 0, 18)
		itemOrder[i] = make([]*ItemOrder, 0, 30)
	}
	// travel timeline
	for _, frame := range matchTL.Info.Frames {
		for _, event := range frame.Events {
			parId := event.ParticipantId - 1
			switch event.Type {
			case "SKILL_LEVEL_UP":
				skillSlot[parId] =
					append(skillSlot[parId], event.SkillSlot)
			case "ITEM_PURCHASED":
				itemOrder[parId] =
					append(itemOrder[parId], &ItemOrder{
						ItemID:    event.ItemId,
						TimeStamp: time.Duration(event.GameTimestamp) * time.Millisecond,
					})
			case "ITEM_SOLD":
				itemOrder[parId] =
					append(itemOrder[parId], &ItemOrder{
						ItemID:    event.ItemId,
						Sold:      true,
						TimeStamp: time.Duration(event.GameTimestamp) * time.Millisecond,
					})
			case "ITEM_UNDO":
				// undo => pop item from ord
				size := len(itemOrder[parId])
				itemOrder[parId] = itemOrder[parId][:size-1]
			}
		}
	}

	var (
		buff []byte
		err  error
	)
	for i := 0; i < partySize; i++ {
		partics[i].Build.SkillOrder = skillSlot[i]
		if buff, err = json.Marshal(skillSlot[i]); err != nil {
			return err
		} else {
			partics[i].Build.Skill = string(buff)
		}

		partics[i].Build.ItemOrder = itemOrder[i]
		if buff, err = json.Marshal(itemOrder[i]); err != nil {
			return err
		} else {
			partics[i].Build.Item = string(buff)
		}
	}
	m.Participants = partics
	return nil
}
