package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/cralack/ChaosMetrics/server/internal/global"
	"github.com/cralack/ChaosMetrics/server/model"
	"github.com/cralack/ChaosMetrics/server/utils"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UsrService struct {
	db     *gorm.DB
	rdb    *redis.Client
	logger *zap.Logger
	rlock  *sync.RWMutex
}

var tokenLen = 32

func NewUserService() *UsrService {
	return &UsrService{
		db:     global.ChaDB,
		rdb:    global.ChaRDB,
		logger: global.ChaLogger,
		rlock:  &sync.RWMutex{},
	}
}

func (s *UsrService) PreRegister(tar *model.User) (string, error) {
	var (
		userDB = &model.User{}
		err    error
	)
	if err = s.db.Where("email=?", tar.Email).First(userDB).Error; err != nil && !userDB.CreatedAt.IsZero() {
		return "", errors.New("email已经被注册，请重试")
	}
	if err = s.db.Where("username=?", tar.UserName).First(userDB).Error; err != nil && !userDB.CreatedAt.IsZero() {
		return "", errors.New("用户名已经被注册，请重试")
	}

	tar.UUID = uuid.Must(uuid.NewRandom())
	tar.Token = utils.GenerateRandomString(tokenLen)
	tar.NickName = "nick" + utils.GenerateRandomString(12)
	tar.Role = model.Civilian
	key := fmt.Sprintf("user:register-%s", tar.Token)
	if err = s.rdb.Set(context.Background(), key, tar, time.Hour*24).Err(); err != nil {
		return "", err
	}
	return tar.Token, nil
}

func (s *UsrService) SendVerifyEmail(token string) error {

	return nil
}

func (s *UsrService) VerifyRegister(token string) (bool, error) {
	var err error
	key := fmt.Sprintf("user:register-%s", token)
	val := s.rdb.Get(context.Background(), key).Val()
	var tar *model.User
	if err = json.Unmarshal([]byte(val), &tar); err != nil {
		return false, err
	}
	if tar.Token != token {
		return false, nil
	}
	if err = s.rdb.Del(context.Background(), key).Err(); err != nil {
		return false, err
	}

	// gen passwd

	tar.Password = utils.BcryptHash(tar.Password)
	if err = s.db.Create(tar).Error; err != nil {
		return false, err
	}

	return true, nil
}

func (s *UsrService) Login(usrname, passwd string) (*model.User, error) {
	var tarDB *model.User
	if err := s.db.Where("username=?", usrname).First(&tarDB).Error; err != nil && tarDB.NickName == "" {
		return nil, errors.New("user does not exist")
	}
	if ok := utils.BcryptCheck(passwd, tarDB.Password); !ok {
		return nil, errors.New("wrong password")
	}
	return tarDB, nil
}

func (s *UsrService) GetUserIno(uuid uuid.UUID) (*model.User, error) {
	tar := &model.User{}
	if err := s.db.Where("uuid=?", uuid).First(&tar).Error; err != nil {
		return nil, err
	}
	tar.ID = 0
	tar.Password = "******"
	return tar, nil
}

func (s *UsrService) ChangePassword(id uint, oldPasswd, newPasswd string) (err error) {
	var tar *model.User
	if err = global.ChaDB.Where("id=?", id).First(&tar).Error; err != nil {
		return err
	}

	if ok := utils.BcryptCheck(oldPasswd, tar.Password); !ok {
		return errors.New("passwd check failed")
	}
	tar.Password = utils.BcryptHash(newPasswd)
	err = global.ChaDB.Save(&tar).Error
	return err
}
