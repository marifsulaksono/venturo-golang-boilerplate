package structs

import "time"

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
		UserID       string `json:"user_id" gorm:"not null;type:char(36);unique"`
		RefreshToken string `json:"refresh_token" gorm:"not null"`
		IP           string `json:"ip" gorm:"not null;type:char(100)"`
	}

	RefreshAccessToken struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}
)
