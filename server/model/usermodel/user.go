package usermodel

import (
	"encoding/json"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	UUID     uuid.UUID `json:"uuid" gorm:"index"`
	UserName string    `gorm:"column:username"`
	Password string    `gorm:"column:password"`
	Email    string    `gorm:"column:email"`
	NickName string    `gorm:"column:nickname"`
	Role     uint      `gorm:"column:role"`

	Token string `gorm:"-"` // 注册token或者登录token
}

func (u *User) MarshalBinary() ([]byte, error) { return json.Marshal(u) }

func (u *User) UnmarshalBinary(bt []byte) error { return json.Unmarshal(bt, u) }
