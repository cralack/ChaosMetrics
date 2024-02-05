package utils

import (
	"github.com/cralack/ChaosMetrics/server/internal/global"
	"golang.org/x/crypto/bcrypt"
)

// BcryptHash 使用 bcrypt 对密码进行加密
func BcryptHash(password string) string {
	var cost int
	switch global.ChaEnv {
	case global.TestEnv:
		cost = bcrypt.MinCost
	case global.DevEnv:
		cost = bcrypt.DefaultCost
	case global.ProductEnv:
		cost = bcrypt.MaxCost
	}
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), cost)
	return string(bytes)
}

// BcryptCheck 对比明文密码和数据库的哈希值
func BcryptCheck(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
