package test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"
	
	"github.com/cralack/ChaosMetrics/server/global"
	"github.com/cralack/ChaosMetrics/server/model/riotmodel"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type User struct {
	gorm.Model
	Name string `gorm:"primary_key;column:user_name;type:varchar(100)"`
	Sex  bool
	Age  int
}

func Test_db_crud_func(t *testing.T) {
	db := global.GVA_DB
	if err := db.AutoMigrate(&User{}); err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 1)
	// 清空表中的所有数据
	if err := db.Exec("truncate table users").Error; err != nil {
		t.Fatal(err)
	}
	
	// 增
	db.Create(&User{
		Name: "snoop",
		Sex:  false,
		Age:  18,
	})
	db.Create(&User{
		Name: "grant",
		Sex:  true,
		Age:  20,
	})
	db.Create(&User{
		Name: "rui",
		Sex:  true,
		Age:  22,
	})

	// 查
	var tar User
	db.First(&tar, "user_name=?", "grant")
	fmt.Println("first:\n", tar.Name, tar.ID)

	var tars []User
	// snop+grant		+rui
	db.Where("age<?", 21).Or("id>?", 2).Find(&tars)
	fmt.Println("find:")
	for _, t := range tars {
		fmt.Println(t.ID, t.Name, t.Age)
	}

	// 改
	db.Where("id=?", 1).First(&User{}).Update("user_name", "snoop dogg")
	db.Where("id in (?)", []int{1, 2}).Find(&[]User{}).Updates(
		map[string]interface{}{
			"Name": "after",
			"Sex":  true,
			"Age":  19,
		})
	
	// 删
	db.Where("id in (?)", []int{1, 3}).Delete(&User{})
	db.Where("id=?", 2).Unscoped().Delete(&User{})
}

func Test_match_summoners(t *testing.T) {
	/*
		set up values
	*/
	// load local json data
	buff, err := os.ReadFile(path + "summoners.txt")
	if err != nil {
		t.Fatal(err)
	}
	// parse
	var summoners []*riotmodel.SummonerDTO
	err = json.Unmarshal(buff, &summoners)
	if err != nil {
		t.Fatal(err)
	}
	// load local json data
	buff, err = os.ReadFile(path + "match.txt")
	if err != nil {
		t.Fatal(err)
	}
	// parse
	var match *riotmodel.MatchDto
	err = json.Unmarshal(buff, &match)
	if err != nil {
		t.Fatal(err)
	}
	
	/*
		test
	*/
	for _, sum := range summoners {
		sum.Matches = append(sum.Matches, match)
		match.Summoners = append(match.Summoners, sum)
	}
	
	db := global.GVA_DB
	db.Exec("DROP TABLE IF EXISTS match_dtos, match_summoners, participant_dtos, summoner_dtos, team_dtos")
	// AutoMigrate
	if err := db.AutoMigrate(
		&riotmodel.MatchDto{},
		&riotmodel.ParticipantDto{},
		&riotmodel.TeamDto{},
		&riotmodel.SummonerDTO{},
	); err != nil {
		t.Fatal(err)
	}
	
	// save data
	if err := db.Save(match).Error; err != nil {
		t.Log(err)
	}
	t.Log("all model store succeed")
	
	// load data 1
	var tar1 *riotmodel.MatchDto
	if err = db.Where("meta_match_id", "TW2_81882122").Preload(
		clause.Associations).First(&tar1).Error; err != nil {
		t.Log(err)
	}
	t.Log(tar1.Metadata.DataVersion)
	
	// load data 2
	tar2 := &riotmodel.MatchDto{
		Metadata: &riotmodel.MetadataDto{
			MetaMatchID: "TW2_81882122",
		},
	}
	if err = db.Preload(clause.Associations).First(&tar2).Error; err != nil {
		t.Log(err)
	}
	
	/*
		Currently it looks like the model doesn't need to change or delete data
	*/
}

// // Minimatch 表示一个用于展示比赛关系的模型。
// // 它包含了 Match 结构体中的部分属性，提供了一个简化的视图。
// // Minimatch represents a model for displaying the relationship between matches.
// // It contains a subset of properties from the Match struct, providing a simplified view.
// type MiniMatchDto struct {
// 	gorm.Model
// 	Summoners []*MiniSummonerDTO `gorm:"many2many:match_summoners"`
//
// 	Metadata *MetadataDto `json:"metadata" gorm:"embedded"` // 比赛元数据
// 	Info     *InfoDto     `json:"info" gorm:"embedded"`     // 比赛信息
// }
// type MetadataDto struct {
// 	DataVersion  string   `json:"dataVersion" gorm:"column:data_version" ` // 比赛数据版本
// 	MatchID      string   `json:"matchId" gorm:"index;column:match_id"`    // 比赛ID
// 	Participants []string `json:"participants" gorm:"-"`                   // 参与者 PUUID 列表
// }
//
// type InfoDto struct {
// 	GameMode     string                `json:"gameMode" gorm:"column:game_mode"`        // 游戏模式
// 	GameVersion  string                `json:"gameVersion" gorm:"column:game_version"`  // 游戏版本
// 	Participants []*MiniParticipantDto `json:"participants" gorm:"foreignKey:match_id"` // 参与者列表
// }
//
// func (p *InfoDto) UnmarshalJSON(data []byte) error {
// 	var f map[string]interface{}
// 	err := json.Unmarshal(data, &f)
// 	if err != nil {
// 		return err
// 	}
// 	for k, v := range f {
// 		switch k {
// 		case "gameMode":
// 			p.GameMode = v.(string)
// 		case "gameVersion":
// 			p.GameVersion = v.(string)
// 		case "participants":
// 			participantsJSON, err := json.Marshal(v)
// 			if err != nil {
// 				return err
// 			}
// 			err = json.Unmarshal(participantsJSON, &p.Participants)
// 			if err != nil {
// 				return err
// 			}
// 			for _, par := range p.Participants {
// 				par.Perks = &Perks{}
// 				par.PerksMeta.parsePerksMeta(par.Perks)
// 			}
// 		}
// 	}
//
// 	return nil
// }
// func (p *MetaPerksDto) parsePerksMeta(buff *Perks) {
// 	buff.StatDefence = p.StatPerks.Defense
// 	buff.StatFlex = p.StatPerks.Flex
// 	buff.StatOffense = p.StatPerks.Offense
// 	for _, sty := range p.Styles {
// 		switch sty.Description {
// 		case "primaryStyle":
// 			buff.PrimarySelection0 = sty.Selections[0].Perk
// 			buff.PrimarySelection1 = sty.Selections[1].Perk
// 			buff.PrimarySelection2 = sty.Selections[2].Perk
// 			buff.PrimarySelection3 = sty.Selections[3].Perk
// 		case "subStyle":
// 			buff.SubSelection0 = sty.Selections[0].Perk
// 			buff.SubSelection1 = sty.Selections[1].Perk
// 		}
// 	}
// }
//
// type MiniParticipantDto struct {
// 	gorm.Model
// 	MatchID      string        `gorm:"column:match_id;index"`                    // 比赛ID
// 	ChampionName string        `json:"championName" gorm:"column:champion_name"` // 英雄名称
// 	Assists      int           `json:"assists" gorm:"column:assists"`            // 助攻数
// 	Deaths       int           `json:"deaths" gorm:"column:deaths"`              // 死亡数
// 	Kills        int           `json:"kills" gorm:"column:kills"`                // 击杀数
// 	PerksMeta    *MetaPerksDto `json:"perks" gorm:"-"`                           // 符文
// 	Perks        *Perks        `gorm:"embedded;embeddedPrefix:perks_"`
// }
//
// type MetaPerksDto struct {
// 	StatPerks struct {
// 		Defense int `json:"defense"`
// 		Flex    int `json:"flex"`
// 		Offense int `json:"offense"`
// 	} `json:"statPerks"`
// 	Styles []struct {
// 		Description string `json:"description"`
// 		Selections  []struct {
// 			Perk int `json:"perk"`
// 			Var1 int `json:"var1"`
// 			Var2 int `json:"var2"`
// 			Var3 int `json:"var3"`
// 		} `json:"selections"`
// 		Style int `json:"style"`
// 	} `json:"styles"`
// }
//
// type Perks struct {
// 	StatDefence       int
// 	StatFlex          int
// 	StatOffense       int
// 	PrimarySelection0 int
// 	PrimarySelection1 int
// 	PrimarySelection2 int
// 	PrimarySelection3 int
// 	SubSelection0     int
// 	SubSelection1     int
// }
//
// func Test_mini_riot_match_store(t *testing.T) {
// 	// load local json data
// 	buff, err := os.ReadFile(path + "match_eg.txt")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// parse
// 	var res MiniMatchDto
// 	err = json.Unmarshal(buff, &res)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	db := global.GVA_DB
// 	// AutoMigrate
// 	if err := db.AutoMigrate(
// 		&MiniMatchDto{},
// 		&MiniParticipantDto{},
// 	); err != nil {
// 		t.Fatal(err)
// 	}
// 	// save data
// 	if err := db.Save(&res).Error; err != nil {
// 		t.Log(err)
// 	}
// }
//
// func Test_riot_match_store(t *testing.T) {
// 	// load local json data
// 	buff, err := os.ReadFile(path + "match.txt")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// parse
// 	var res riotmodel.MatchDto
// 	err = json.Unmarshal(buff, &res)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	db := global.GVA_DB
// 	// AutoMigrate
// 	if err := db.AutoMigrate(
// 		&riotmodel.MatchDto{},
// 		&riotmodel.ParticipantDto{},
// 		&riotmodel.TeamDto{},
// 	); err != nil {
// 		t.Fatal(err)
// 	}
// 	// save data
// 	if err := db.Save(&res).Error; err != nil {
// 		t.Log(err)
// 	}
// }
//
// type MiniSummonerDTO struct {
// 	gorm.Model
// 	Matchs []*MiniMatchDto `gorm:"many2many:match_summoners"`
//
// 	Name  string `gorm:"column:name;type:varchar(100)" json:"name"`
// 	PUUID string `gorm:"column:puuid;type:varchar(100)" json:"puuid"`
// }
//
// func Test_summoner_store(t *testing.T) {
// 	// load local json data
// 	buff, err := os.ReadFile(path + "summoners.txt")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// parse
// 	var res []*riotmodel.SummonerDTO
// 	err = json.Unmarshal(buff, &res)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	db := global.GVA_DB
// 	if err := db.AutoMigrate(
// 		&riotmodel.SummonerDTO{},
// 	); err != nil {
// 		t.Fatal(err)
// 	}
// 	// save data
// 	if err := db.Save(&res).Error; err != nil {
// 		t.Log(err)
// 	}
// }

