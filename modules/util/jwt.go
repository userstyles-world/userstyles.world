package util

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	"userstyles.world/modules/config"
)

type JWTTokenBuilder struct {
	*jwt.Token
}

func NewJWT() *JWTTokenBuilder {
	return &JWTTokenBuilder{jwt.New(jwt.SigningMethodHS512)}
}

func (jt *JWTTokenBuilder) SetClaim(name string, value any) *JWTTokenBuilder {
	jt.Claims.(jwt.MapClaims)[name] = value
	return jt
}

func (jt *JWTTokenBuilder) SetExpiration(duration time.Time) *JWTTokenBuilder {
	if !duration.IsZero() {
		jt.Claims.(jwt.MapClaims)["exp"] = duration.Unix()
	}
	return jt
}

func (jt *JWTTokenBuilder) GetSignedString(customKey []byte) (string, error) {
	if customKey == nil {
		customKey = []byte(config.Secrets.SessionTokenKey)
	}
	return jt.SignedString(customKey)
}
