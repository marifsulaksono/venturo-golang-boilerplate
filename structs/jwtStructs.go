package structs

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTUser struct {
	ID              uuid.UUID   `json:"id"`
	Email           string      `json:"email"`
	UpdatedSecurity interface{} `json:"updated_security"`
	Access          string      `json:"access"`
	jwt.RegisteredClaims
}
