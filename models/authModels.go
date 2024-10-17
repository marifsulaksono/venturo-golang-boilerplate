package models

import (
	"context"
	"errors"
	"simple-crud-rnd/structs"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AuthModel struct {
	db *gorm.DB
}

func NewAuthModel(db *gorm.DB) *AuthModel {
	return &AuthModel{
		db: db,
	}
}

func (am *AuthModel) GetByRefreshToken(ctx context.Context, refreshToken string) (structs.TokenAuth, error) {
	token := structs.TokenAuth{}
	err := am.db.Select("user_id", "refresh_token", "ip").Where("refresh_token = ?", refreshToken).
		First(&token).Error
	return token, err
}

func (am *AuthModel) Upsert(ctx context.Context, payload *structs.TokenAuth) error {
	if err := am.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "user_id"}, {Name: "ip"}}, // Specify columns to check for conflict
		DoUpdates: clause.Assignments(map[string]interface{}{
			"refresh_token": payload.RefreshToken, // Update refresh_token on conflict
		}),
	}).Create(&payload).Error; err != nil {
		return err
	}

	return nil
}

func (am *AuthModel) Update(ctx context.Context, payload *structs.TokenAuth) (structs.TokenAuth, error) {
	token := structs.TokenAuth{RefreshToken: payload.RefreshToken}
	res := am.db.Model(&token).Updates(&payload)
	if res.RowsAffected == 0 {
		return token, errors.New("no rows updated")
	}
	return token, nil
}

func (am *AuthModel) Delete(ctx context.Context, payload *structs.TokenAuth) error {
	return am.db.Where("user_id = ? AND ip = ?", payload.UserID, payload.IP).Delete(&structs.Role{}).Error
}
