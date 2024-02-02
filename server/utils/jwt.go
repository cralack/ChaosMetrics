package utils

import (
	"errors"
	"time"

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

const SigningKey = "Open Sesame"

func NewJWT() *JWT {
	return &JWT{SigningKey: []byte(SigningKey)}
}

func (j *JWT) CreateClaims(baseClaims request.BaseClaims) request.CustomClaims {
	buf := time.Hour * 24
	exp := time.Hour * 2
	return request.CustomClaims{
		BaseClaims: baseClaims,
		BufferTime: int64(buf / time.Second),
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "ChaosMetric",
			Audience:  jwt.ClaimStrings{"Chao"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1000)),
			NotBefore: jwt.NewNumericDate(time.Now().Add(exp)),
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
		if ve, ok := err.(*jwt.ValidationError); ok {
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
	}
	if token != nil {
		if claims, ok := token.Claims.(*request.CustomClaims); ok && token.Valid {
			return claims, nil
		}
	}
	return nil, TokenInvalid
}
