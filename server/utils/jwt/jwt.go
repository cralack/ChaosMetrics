package jwt

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cralack/ChaosMetrics/server/internal/global"
	"github.com/cralack/ChaosMetrics/server/model/request"
	"github.com/golang-jwt/jwt/v4"
)

var (
	TokenExpired     = errors.New("Token is expired")
	TokenNotValidYet = errors.New("Token not active yet")
	TokenMalformed   = errors.New("That's not even a token")
	TokenInvalid     = errors.New("Couldn't handle this token:")
)

type JWT struct {
	SigningKey []byte
}

func NewJWT() *JWT {
	return &JWT{SigningKey: []byte(global.ChaConf.JwtConf.SigningKey)}
}

func (j *JWT) CreateClaims(baseClaims request.PrivateClaims) request.CustomClaims {
	// buf := global.ChaConf.JwtConf.BufferTime
	exp := global.ChaConf.JwtConf.ExpiresTime
	return request.CustomClaims{
		PrivateClaims: baseClaims,
		// BufferTime:    int64(buf / time.Second),
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    global.ChaConf.JwtConf.Issuer,
			Audience:  jwt.ClaimStrings{"Chao"},
			NotBefore: jwt.NewNumericDate(time.Now().Add(-1000)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(exp)),
		},
	}
}

func (j *JWT) CreateToken(claims request.CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

func (j *JWT) ParseToken(tokenStr string) (*request.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &request.CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.SigningKey, nil
	})
	if err != nil {
		var ve *jwt.ValidationError
		if errors.As(err, &ve) {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
		// 如果不是ValidationError，直接返回TokenInvalid
		return nil, TokenInvalid
	}
	if token != nil {
		if claims, ok := token.Claims.(*request.CustomClaims); ok && token.Valid {
			return claims, nil
		}
	}
	return nil, TokenInvalid
}

func (j *JWT) SetJWTBlack(jwt string) (err error) {
	// 此处过期时间等于jwt过期时间
	dura := global.ChaConf.JwtConf.ExpiresTime
	err = global.ChaRDB.Set(context.Background(), fmt.Sprintf("jwt-%s", jwt), 1, dura).Err()
	return err
}

func (j *JWT) InBlackList(jwt string) bool {
	_, err := global.ChaRDB.Get(context.Background(), fmt.Sprintf("jwt-%s", jwt)).Result()
	return err == nil
}
