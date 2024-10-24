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

func (jm *JobModel) GetAll(ctx context.Context, limit, offset int) ([]structs.Job, int64, error) {
	jobs := []structs.Job{}
	if err := jm.db.Select("id", "job_name", "payload", "status", "attempts", "created_at", "updated_at").
		Limit(limit).Offset(offset).Find(&jobs).Error; err != nil {
		return nil, 0, err
	}

	var count int64
	if err := jm.db.Table("jobs").Count(&count).Error; err != nil {
		return nil, 0, err
	}

	return jobs, count, nil
}

func (jm *JobModel) GetById(ctx context.Context, id string) (structs.Job, error) {
	job := structs.Job{}
	err := jm.db.Select("id", "job_name", "payload", "status", "attempts", "created_at", "updated_at").Where("id = ?", id).
		First(&job).Error
	return job, err
}

func (jm *JobModel) Create(ctx context.Context, payload *structs.Job) error {
	return jm.db.Create(payload).Error
}
