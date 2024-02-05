package model

import (
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model

	ChampionID string `json:"championID"`
	Version    string `json:"version"`
	Content    string `json:"content,omitempty"`
	AuthorID   uint   `gorm:"column:author_id;not null"`
	Author     *User  `json:"author,omitempty" gorm:"foreignKey:AuthorID"`
}
