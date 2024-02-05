package comment

import (
	"github.com/cralack/ChaosMetrics/server/model"
)

type CommentsDTO struct {
	ID         uint
	ChampionID string
	Version    string
	Author     *Author
	Comment    string
}
type Author struct {
	NickName string
}

func ConvertCommentsDTO(src ...*model.Comment) (des []*CommentsDTO) {
	if len(src) == 0 {
		return
	}
	des = make([]*CommentsDTO, 0, len(src))
	for _, c := range src {
		tmp := CommentsDTO{
			ID:         c.ID,
			ChampionID: c.ChampionID,
			Version:    c.Version,
			Comment:    c.Content,
		}
		if c.Author != nil {
			tmp.Author = &Author{
				NickName: c.Author.NickName,
			}
		}
		des = append(des, &tmp)
	}
	return
}
