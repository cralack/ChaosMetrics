package riotmodel

import (
	"errors"
	"time"
)

type Version []string

type Image struct {
	Full   string `json:"full"`
	Sprite string `json:"sprite"`
	Group  string `json:"group"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
	W      int    `json:"w"`
	H      int    `json:"h"`
}

func convertTime(v interface{}, layout string) (t time.Time, err error) {
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
