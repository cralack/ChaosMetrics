package config

import (
	"time"
)

type JwtConfig struct {
	SigningKey  string        `mapstructure:"signing_key" yaml:"signing_key"`
	ExpiresTime time.Duration `mapstructure:"expires_time" yaml:"expires_time"`
	BufferTime  time.Duration `mapstructure:"buffer_time" yaml:"buffer_time"`
	Issuer      string        `mapstructure:"issuer" yaml:"issuer"`
}
