package structs

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (Role) TableName() string {
	return "m_role"
}

type (
	Role struct {
		ID        uuid.UUID       `json:"id" gorm:"primaryKey;type:char(36);not null"`
		Name      string          `json:"name" gorm:"type:char(36);unique;not null"`
		Access    string          `json:"access" gorm:"not null"`
		CreatedAt time.Time       `json:"created_at" gorm:"autoCreateTime"`
		UpdatedAt time.Time       `json:"updated_at" gorm:"autoUpdateTime"`
		DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
		CreatedBy *uuid.UUID      `json:"created_by,omitempty" gorm:"type:char(36)"`
		UpdatedBy *uuid.UUID      `json:"updated_by,omitempty" gorm:"type:char(36)"`
		DeletedBy *uuid.UUID      `json:"deleted_by,omitempty" gorm:"type:char(36)"`
	}

	RoleRequest struct {
		ID     uuid.UUID   `json:"id"`
		Name   string      `json:"name" validate:"required"`
		Access interface{} `json:"access"`
	}
)

func (r *Role) BeforeCreate(tx *gorm.DB) error {
	r.ID = uuid.New()
	return nil
}
