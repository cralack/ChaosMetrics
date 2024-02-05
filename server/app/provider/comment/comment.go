package comment

import (
	"errors"
	"fmt"
	"sync"

	"github.com/cralack/ChaosMetrics/server/internal/global"
	"github.com/cralack/ChaosMetrics/server/model"
	"github.com/cralack/ChaosMetrics/server/model/request"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ComntService struct {
	db    *gorm.DB
	rlock *sync.RWMutex
}

func NewCommentService() *ComntService {
	return &ComntService{
		db:    global.ChaDB,
		rlock: &sync.RWMutex{},
	}
}

func (s *ComntService) GetComments(championid, version string, pager *request.Pager) (res []*model.Comment, err error) {
	res = make([]*model.Comment, 0, pager.Size)
	var total int64

	if err = s.db.Model(&model.Comment{}).Where("champion_id=?", championid).Where("version=?",
		version).Count(&total).Error; err == nil {
		pager.Total = total
	}

	if err = s.db.Where("champion_id=?", championid).Where("version=?", version).Preload(clause.Associations).Order("created_at desc").Offset(pager.Start).Limit(pager.Size).Find(&res).Error; err != nil {
		if len(res) == 0 {
			return []*model.Comment{}, errors.New(fmt.Sprintf("cannot find comments by %s or %s", championid, version))
		}
	}
	return res, nil
}

func (s *ComntService) GetSingleComment(id uint) (res *model.Comment, err error) {
	if err = s.db.First(res, id).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func (s *ComntService) PostComment(tar *model.Comment) error {
	if tar.AuthorID == 0 {
		return errors.New("author ID cannot be empty")
	}
	if err := s.db.Create(tar).Error; err != nil {
		return err
	}
	return nil
}

func (s *ComntService) DeleteComments(id uint) error {
	result := s.db.Where("id=?", id).Delete(&model.Comment{})
	if err := result.Error; err != nil {
		return err
	}
	if result.RowsAffected == 0 {
		return errors.New("no record delete")
	}
	return nil
}
