package test

import (
	"testing"

	"github.com/cralack/ChaosMetrics/server/model"
	"gorm.io/gorm/clause"
)

func Test_comment(t *testing.T) {
	coment := &model.Comment{
		ChampionID: "Ahri",
		Version:    "14.2.1",
		Content:    "so cute",
		AuthorID:   2,
	}
	if err := db.Create(coment).Error; err != nil {
		t.Log(err)
	}
	var tar *model.Comment
	if err := db.Where("champion_id=?", "Ahri").Preload(
		clause.Associations).Find(&tar).Error; err != nil {
		t.Log(err)
	}
}

func Test_comment_list(t *testing.T) {
	championid := "Ahri"
	version := "14.1.1"
	var total int64
	db.Model(&model.Comment{}).Where("champion_id = ?", championid).Where("version = ?", version).Count(&total)
	t.Log(total)
}
