package models

import (
	"context"
	"errors"
	"simple-crud-rnd/helpers"
	"simple-crud-rnd/structs"

	"gorm.io/gorm"
)

type AuthModel struct {
	db *gorm.DB
}

func NewAuthModel(db *gorm.DB) *AuthModel {
	return &AuthModel{
		db: db,
	}
}

func (am *AuthModel) GetByUserIDAndIP(ctx context.Context, userId, ip string) (structs.TokenAuth, error) {
	token := structs.TokenAuth{}
	err := am.db.Select("user_id", "refresh_token", "ip").Where("user_id = ? AND ip = ?", userId, ip).
		First(&token).Error
	return token, helpers.SendTraceErrorToSentry(err)
}

func (am *AuthModel) GetByRefreshToken(ctx context.Context, refreshToken string) (structs.TokenAuth, error) {
	token := structs.TokenAuth{}
	err := am.db.Select("user_id", "refresh_token", "ip").Where("refresh_token = ?", refreshToken).
		First(&token).Error
	return token, helpers.SendTraceErrorToSentry(err)
}

func (am *AuthModel) Upsert(ctx context.Context, payload *structs.TokenAuth) error {
	token, err := am.GetByUserIDAndIP(ctx, payload.UserID, payload.IP)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create a new data if not found
			if err := am.db.Create(payload).Error; err != nil {
				return helpers.SendTraceErrorToSentry(err)
			}
			return nil
		}
		return helpers.SendTraceErrorToSentry(err)
	}

	// update if exists user_id and ip
	if err := am.db.Model(&structs.TokenAuth{}).
		Where("user_id = ? AND ip = ?", token.UserID, token.IP).
		Update("refresh_token", payload.RefreshToken).Error; err != nil {
		return helpers.SendTraceErrorToSentry(err)
	}

	return nil
}

func (am *AuthModel) Delete(ctx context.Context, refreshToken, ip string) error {
	err := am.db.Where("refresh_token = ? AND ip = ?", refreshToken, ip).Delete(&structs.TokenAuth{}).Error
	return helpers.SendTraceErrorToSentry(err)
}
