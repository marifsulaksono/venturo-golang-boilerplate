package models

import (
	"context"
	"simple-crud-rnd/structs"

	"gorm.io/gorm"
)

type JobModel struct {
	db *gorm.DB
}

func NewJobModel(db *gorm.DB) *JobModel {
	return &JobModel{
		db: db,
	}
}

func (jm *JobModel) Create(ctx context.Context, payload *structs.Job) error {
	return jm.db.Create(payload).Error
}
