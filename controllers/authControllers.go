package controllers

import (
	"errors"
	"net/http"
	"simple-crud-rnd/config"
	"simple-crud-rnd/helpers"
	"simple-crud-rnd/models"
	"simple-crud-rnd/structs"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthController struct {
	db        *gorm.DB
	authModel *models.AuthModel
	userModel *models.UserModel
	cfg       *config.Config
}

func NewAuthController(db *gorm.DB, authModel *models.AuthModel, userModel *models.UserModel, cfg *config.Config) *AuthController {
	return &AuthController{db, authModel, userModel, cfg}
}

func (ah *AuthController) Login(c echo.Context) error {
	var (
		ctx     = c.Request().Context()
		ip      = c.RealIP()
		request structs.Login
	)

	if err := c.Bind(&request); err != nil {
		return helpers.Response(c, http.StatusBadRequest, nil, err.Error())
	}

	if err := c.Validate(&request); err != nil {
		return helpers.Response(c, http.StatusBadRequest, nil, err.Error())
	}

	user, err := ah.userModel.GetByEmail(ctx, request.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helpers.Response(c, http.StatusNotFound, nil, "Email atau password salah")
		}
		return helpers.Response(c, http.StatusInternalServerError, nil, err.Error())
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		return helpers.Response(c, http.StatusNotFound, nil, "Email atau password salah")
	}

	accessToken, expiredAt, err := helpers.GenerateTokenJWT(&user, false)
	if err != nil {
		return helpers.Response(c, http.StatusInternalServerError, nil, err.Error())
	}

	refreshToken, _, err := helpers.GenerateTokenJWT(&user, true)
	if err != nil {
		return helpers.Response(c, http.StatusInternalServerError, nil, err.Error())
	}

	if err := ah.authModel.Upsert(ctx, &structs.TokenAuth{
		UserID:       user.ID.String(),
		RefreshToken: refreshToken,
		IP:           ip,
	}); err != nil {
		return helpers.Response(c, http.StatusInternalServerError, nil, err.Error())
	}

	response := structs.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    *expiredAt,
		Metadata: structs.Metadata{
			Name:   user.Name,
			Email:  user.Email,
			Access: user.Role.Access,
		},
	}

	return helpers.Response(c, http.StatusOK, response, "Login berhasil")
}

func (ah *AuthController) RefreshAccessToken(c echo.Context) error {
	var (
		ctx     = c.Request().Context()
		request structs.RefreshAccessToken
	)

	if err := c.Bind(&request); err != nil {
		return helpers.Response(c, http.StatusBadRequest, nil, err.Error())
	}

	if err := c.Validate(&request); err != nil {
		return helpers.Response(c, http.StatusBadRequest, nil, err.Error())
	}

	token, err := ah.authModel.GetByRefreshToken(ctx, request.RefreshToken)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helpers.Response(c, http.StatusUnauthorized, nil, "Data token tidak ada, harap melakukan login")
		}
		return helpers.Response(c, http.StatusInternalServerError, nil, err.Error())
	}

	user, err := helpers.VerifyTokenJWT(request.RefreshToken, true)
	if err != nil {
		return helpers.Response(c, http.StatusUnauthorized, nil, err.Error())
	}

	accessToken, expiredAt, err := helpers.GenerateTokenJWT(user, false)
	if err != nil {
		return helpers.Response(c, http.StatusInternalServerError, nil, err.Error())
	}

	response := structs.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: token.RefreshToken,
		ExpiresAt:    *expiredAt,
		Metadata: structs.Metadata{
			Name:   user.Name,
			Email:  user.Email,
			Access: user.Role.Access,
		},
	}

	return helpers.Response(c, http.StatusOK, response, "Berhasil")
}

func (ah *AuthController) Logout(c echo.Context) error {
	var (
		ctx     = c.Request().Context()
		ip      = c.RealIP()
		request structs.RefreshAccessToken
	)

	if err := c.Bind(&request); err != nil {
		return helpers.Response(c, http.StatusBadRequest, nil, err.Error())
	}

	if err := c.Validate(&request); err != nil {
		return helpers.Response(c, http.StatusBadRequest, nil, err.Error())
	}

	err := ah.authModel.Delete(ctx, request.RefreshToken, ip)
	if err != nil {
		return helpers.Response(c, http.StatusInternalServerError, nil, err.Error())
	}

	return helpers.Response(c, http.StatusOK, nil, "Logout berhasil")
}

func (ah *AuthController) Generate2FASecretKey(c echo.Context) error {
	var (
		ctx = c.Request().Context()
	)

	userId := c.Get("user_id").(uuid.UUID)
	if userId == uuid.Nil {
		return helpers.Response(c, http.StatusUnauthorized, nil, "User ID not found in context")
	}

	user, err := ah.userModel.GetById(ctx, userId)
	if err != nil {
		return helpers.Response(c, http.StatusInternalServerError, nil, err.Error())
	}

	key, err := helpers.Generate2FASecretKey(ctx)
	if err != nil {
		return helpers.Response(c, http.StatusInternalServerError, nil, err.Error())
	}

	if err := ah.authModel.Store2FASecretKey(ctx, &structs.TOTP{UserID: user.ID.String(), SecretKey: key}); err != nil {
		return helpers.Response(c, http.StatusInternalServerError, nil, err.Error())
	}

	qrUrl := helpers.Generate2FAQRCodeURL(user.ID.String(), key)
	return helpers.Response(c, http.StatusOK, map[string]string{"qrcode_url": qrUrl}, "Berhasil")
}

func (ah *AuthController) Verify2FA(c echo.Context) error {
	var (
		request structs.TOTPRequest
	)

	userId := c.Get("user_id").(uuid.UUID)
	if userId == uuid.Nil {
		return helpers.Response(c, http.StatusUnauthorized, nil, "User ID not found in context")
	}

	if err := c.Bind(&request); err != nil {
		return helpers.Response(c, http.StatusBadRequest, nil, err.Error())
	}

	if err := c.Validate(&request); err != nil {
		return helpers.Response(c, http.StatusBadRequest, nil, err.Error())
	}

	totp, err := ah.authModel.Get2FASecretKeyByUserID(c.Request().Context(), userId.String())
	if err != nil {
		return helpers.Response(c, http.StatusInternalServerError, nil, err.Error())
	}

	ok, err := helpers.VerifyOTP(totp.SecretKey, request.Code)
	if err != nil {
		return helpers.Response(c, http.StatusInternalServerError, nil, err.Error())
	} else if !ok {
		return helpers.Response(c, http.StatusBadRequest, nil, "Kode salah")
	}

	return helpers.Response(c, http.StatusOK, nil, "Berhasil")
}
