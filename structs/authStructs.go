package structs

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (TokenAuth) TableName() string {
	return "m_token_auth"
}

type (
	Login struct {
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	LoginResponse struct {
		AccessToken  string    `json:"access_token"`
		RefreshToken string    `json:"refresh_token"`
		ExpiresAt    time.Time `json:"expired_at"`
		Metadata     Metadata  `json:"metadata"`
	}

	Metadata struct {
		Name   string `json:"name"`
		Email  string `json:"email"`
		Access string `json:"Access"`
	}

	TokenAuth struct {
		ID           uuid.UUID `json:"id" gorm:"primaryKey;type:char(36);not null"`
		UserID       string    `json:"user_id" gorm:"not null;type:char(36)"`
		RefreshToken string    `json:"refresh_token" gorm:"not null"`
		IP           string    `json:"ip" gorm:"not null;type:char(100)"`
	}

	RefreshAccessToken struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}
)

func (u *TokenAuth) BeforeCreate(tx *gorm.DB) error {
	u.ID = uuid.New()
	return nil
}
