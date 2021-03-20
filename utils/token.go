package utils

import (
	"time"

	"github.com/form3tech-oss/jwt-go"
	"userstyles.world/config"
)

type JWTTokenBuilder struct {
	*jwt.Token
}

func NewJWTToken() *JWTTokenBuilder {
	return &JWTTokenBuilder{jwt.New(jwt.SigningMethodHS512)}
}

func (jt *JWTTokenBuilder) SetClaim(name string, value interface{}) *JWTTokenBuilder {
	jt.Claims.(jwt.MapClaims)[name] = value
	return jt
}

func (jt *JWTTokenBuilder) SetExpiration(duration time.Duration) *JWTTokenBuilder {
	jt.Claims.(jwt.MapClaims)["exp"] = time.Now().Add(duration).Unix()
	return jt
}

func (jt *JWTTokenBuilder) GetSignedString() (string, error) {
	return jt.SignedString([]byte(config.JWT_SIGNING_KEY))
}
