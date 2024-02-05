package request

import (
	"github.com/cralack/ChaosMetrics/server/model"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type CustomClaims struct {
	PrivateClaims
	// BufferTime int64
	jwt.RegisteredClaims
}

type PrivateClaims struct {
	UUID     uuid.UUID
	ID       uint
	Username string
	NickName string
	Role     model.Role
}
