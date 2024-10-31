package models

import (
	"context"
	"errors"
	"simple-crud-rnd/helpers"
	"simple-crud-rnd/structs"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RoleModel struct {
	db *gorm.DB
}

func NewRoleModel(db *gorm.DB) *RoleModel {
	return &RoleModel{
		db: db,
	}
}

func (rm *RoleModel) GetAll(ctx context.Context, limit, offset int) ([]structs.Role, int64, error) {
	roles := []structs.Role{}
	if err := rm.db.Select("id", "name", "access", "created_at", "updated_at").
		Limit(limit).Offset(offset).Find(&roles).Error; err != nil {
		return nil, 0, helpers.SendTraceErrorToSentry(err)
	}

	var count int64
	if err := rm.db.Table("m_role").Where("deleted_at IS NULL").Count(&count).Error; err != nil {
		return nil, 0, helpers.SendTraceErrorToSentry(err)
	}

	return roles, count, nil
}

func (rm *RoleModel) GetById(ctx context.Context, id uuid.UUID) (structs.Role, error) {
	role := structs.Role{}
	err := rm.db.Select("id", "name", "access", "created_at", "updated_at", "created_by", "updated_by").
		Where("deleted_at IS NULL").First(&role, id).Error
	return role, helpers.SendTraceErrorToSentry(err)
}

func (rm *RoleModel) Create(ctx context.Context, payload *structs.Role) (structs.Role, error) {
	role := structs.Role{
		Name:   payload.Name,
		Access: payload.Access,
	}

	res := rm.db.Create(&role).Clauses(clause.Returning{
		Columns: []clause.Column{
			{Name: "id"},
			{Name: "name"},
			{Name: "access"},
		},
	})

	if res.Error != nil {
		return role, helpers.SendTraceErrorToSentry(res.Error)
	}

	return role, nil
}

func (rm *RoleModel) Update(ctx context.Context, payload *structs.Role) (structs.Role, error) {
	role := structs.Role{ID: payload.ID}
	res := rm.db.Model(&role).Clauses(clause.Returning{}).Updates(&payload)
	if res.RowsAffected == 0 {
		return role, helpers.SendTraceErrorToSentry(errors.New("no rows updated"))
	}
	return role, nil
}

func (rm *RoleModel) Delete(ctx context.Context, id uuid.UUID) error {
	res := rm.db.Delete(&structs.Role{}, id)
	if res.RowsAffected == 0 {
		return helpers.SendTraceErrorToSentry(errors.New("no rows deleted"))
	}
	return nil
}
