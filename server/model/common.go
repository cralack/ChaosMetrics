package model

import (
	"errors"
	"time"
)

// type Storable interface {
// 	TableName() string
// }

type Image struct {
	Full   string `json:"full" gorm:"column:full"`     // 大图文件名
	Sprite string `json:"sprite" gorm:"column:sprite"` // 小图文件名
	Group  string `json:"group" gorm:"column:group"`   // 图像所属组
	X      int    `json:"x" gorm:"column:x"`           // 图像 X 坐标
	Y      int    `json:"y" gorm:"column:y"`           // 图像 Y 坐标
	W      int    `json:"w" gorm:"column:w"`           // 图像宽度
	H      int    `json:"h" gorm:"column:h"`           // 图像高度
}

func ConvertTime(v interface{}, layout string) (t time.Time, err error) {
	if timestmp, ok := v.(string); !ok {
		return t, errors.New("val wrong type")
	} else {
		if t, err = time.Parse(layout, timestmp); err != nil {
			return t, err
		} else {
			return t, nil
		}
	}
}
