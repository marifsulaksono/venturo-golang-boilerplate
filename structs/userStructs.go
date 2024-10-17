package structs

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (User) TableName() string {
	return "m_user"
}

type (
	User struct {
		ID              uuid.UUID       `json:"id" gorm:"primaryKey;type:char(36);not null"`
		CreatedAt       time.Time       `json:"created_at" gorm:"autoCreateTime"`
		UpdatedAt       time.Time       `json:"updated_at" gorm:"autoUpdateTime"`
		DeletedAt       *gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
		CreatedBy       *uuid.UUID      `json:"created_by,omitempty" gorm:"type:char(36)"`
		UpdatedBy       *uuid.UUID      `json:"updated_by,omitempty" gorm:"type:char(36)"`
		DeletedBy       *uuid.UUID      `json:"deleted_by,omitempty" gorm:"type:char(36)"`
		Name            string          `json:"name" gorm:"not null"`
		Email           string          `json:"email" gorm:"unique;not null"`
		Photo           string          `json:"photo_url,omitempty"`
		PhoneNumber     string          `json:"phone_number" gorm:"not null"`
		Password        string          `json:"password,omitempty" gorm:"not null"`
		UserRolesId     string          `json:"user_roles_id" gorm:"type:char(36)"`
		Role            *Role           `gorm:"-:migration;foreignKey:UserRolesId;references:ID" json:"role_detail,omitempty"`
		UpdatedSecurity time.Time       `json:"updated_security"`
	}

	UserRequest struct {
		Name            string    `json:"name" validate:"required"`
		Email           string    `json:"email" validate:"required,email"`
		Photo           string    `json:"photo_url,omitempty"`
		PhoneNumber     string    `json:"phone_number" validate:"required,e164"`
		Password        string    `json:"password" validate:"required"`
		UserRolesId     string    `json:"user_roles_id" validate:"required,uuid"`
		UpdatedSecurity time.Time `json:"updated_security"`
	}
)

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.ID = uuid.New()
	return nil
}
