package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/cralack/ChaosMetrics/server/internal/global"
	"github.com/cralack/ChaosMetrics/server/model/usermodel"
	"github.com/cralack/ChaosMetrics/server/utils"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
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

func (s *UsrService) PreRegister(tar *usermodel.User) (string, error) {
	var (
		userDB = &usermodel.User{}
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
	key := fmt.Sprintf("user:register-%s", token)
	val := s.rdb.Get(context.Background(), key).Val()
	var tar *usermodel.User
	if err := json.Unmarshal([]byte(val), &tar); err != nil {
		return false, err
	}
	if tar.Token != token {
		return false, nil
	}
	if err := s.rdb.Del(context.Background(), key).Err(); err != nil {
		return false, err
	}

	// gen passwd
	var cost int
	switch global.ChaEnv {
	case global.TestEnv:
		cost = bcrypt.MinCost
	case global.DevEnv:
		cost = bcrypt.DefaultCost
	case global.ProductEnv:
		cost = bcrypt.MaxCost
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(tar.Password), cost)
	if err != nil {
		return false, err
	}
	tar.Password = string(hash)
	if err = s.db.Create(tar).Error; err != nil {
		return false, err
	}

	return true, nil
}

func (s *UsrService) Login(usrname, passwd string) (string, error) {
	var tarDB *usermodel.User
	if err := s.db.Where("username=?", usrname).First(&tarDB).Error; err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(tarDB.Password), []byte(passwd)); err != nil {
		return "", errors.New("wrong password")
	}
	// cache session
	tarDB.Password = ""
	tarDB.Token = utils.GenerateRandomString(tokenLen)
	key := fmt.Sprintf("session:%s", tarDB.Token)
	if err := s.rdb.Set(context.Background(), key, tarDB, time.Hour*24).Err(); err != nil {
		return "", err
	}
	return tarDB.Token, nil
}

func (s *UsrService) Logout(tar *usermodel.User) error {
	session, err := s.VerifyLogin(tar.Token)
	if err != nil || session.UserName != tar.UserName {
		return errors.New("verify user failed " + err.Error())
	}
	key := fmt.Sprintf("session:%s", tar.Token)
	if err = s.rdb.Del(context.Background(), key).Err(); err != nil {
		return err
	}
	return nil
}

func (s *UsrService) VerifyLogin(token string) (*usermodel.User, error) {
	key := fmt.Sprintf("session:%s", token)
	res := s.rdb.Get(context.Background(), key)
	if res.Err() != nil {
		return nil, res.Err()
	}

	var tar *usermodel.User
	err := json.Unmarshal([]byte(res.Val()), &tar)
	if err != nil {
		return nil, err
	}
	return tar, nil
}

func (s *UsrService) GetUser(userID uint) (*usermodel.User, error) {
	tar := &usermodel.User{}
	tar.ID = userID
	if err := s.db.Where("id=?", userID).First(&tar).Error; err != nil {
		return nil, err
	}
	return tar, nil
}
