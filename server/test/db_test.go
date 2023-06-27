package test

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/cralack/ChaosMetrics/server/global"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name string `gorm:"primary_key;column:user_name;type:varchar(100)"`
	Sex  bool
	Age  int
}

func Test_db(t *testing.T) {
	db := global.GVA_DB
	db.AutoMigrate(&User{})
	time.Sleep(time.Second * 1)
	//清空表中的所有数据
	if err := db.Exec("truncate table users").Error; err != nil {
		log.Fatal(err)
	}

	//增
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

	//查
	var tar User
	db.First(&tar, "user_name=?", "grant")
	fmt.Println("first:\n", tar.Name, tar.ID)

	var tars []User
	//snop+grant		+rui
	db.Where("age<?", 21).Or("id>?", 2).Find(&tars)
	fmt.Println("find:")
	for _, t := range tars {
		fmt.Println(t.ID, t.Name, t.Age)
	}

	//改
	db.Where("id=?", 1).First(&User{}).Update("user_name", "snoop dogg")
	db.Where("id in (?)", []int{1, 2}).Find(&[]User{}).Updates(
		map[string]interface{}{
			"Name": "after",
			"Sex":  true,
			"Age":  19,
		})

	//删
	db.Where("id in (?)", []int{1, 3}).Delete(&User{})
	db.Where("id=?", 2).Unscoped().Delete(&User{})
}
