package structs

import (
	"time"

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

type (
	Login struct {
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	LoginResponse struct {
		AccessToken string    `json:"access_token"`
		ExpiresAt   time.Time `json:"expired_at"`
		Metadata    Metadata  `json:"metadata"`
	}

	Metadata struct {
		Name   string `json:"name"`
		Email  string `json:"email"`
		Access string `json:"Access"`
	}
)
