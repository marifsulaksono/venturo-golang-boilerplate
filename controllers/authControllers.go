package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"simple-crud-rnd/config"
	"simple-crud-rnd/helpers"
	"simple-crud-rnd/models"
	"simple-crud-rnd/rabbitmq"
	"simple-crud-rnd/structs"

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

func (ah *AuthController) ForgotPassword(c echo.Context) error {
	var (
		ctx  = c.Request().Context()
		body map[string]interface{}
	)

	if err := json.NewDecoder(c.Request().Body).Decode(&body); err != nil {
		return helpers.Response(c, http.StatusBadRequest, nil, err.Error())
	}

	email := body["email"].(string)
	user, err := ah.userModel.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helpers.Response(c, http.StatusNotFound, nil, "Email tidak ditemukan")
		}
		return helpers.Response(c, http.StatusInternalServerError, nil, err.Error())
	}

	message, err := helpers.EncryptMessageRSA(user.ID.String())
	if err != nil {
		return helpers.Response(c, http.StatusInternalServerError, nil, err.Error())
	}

	// send email using sendgrid
	/*
		err = helpers.SendMailSendgrid(message, "Lupa Password", user.Name, user.Email)
		if err != nil {
			return helpers.Response(c, http.StatusInternalServerError, nil, err.Error())
		}
	*/

	// send email using gmail
	err = helpers.SendMailGmail(message, "Lupa Password", user.Email)
	if err != nil {
		return helpers.Response(c, http.StatusInternalServerError, nil, err.Error())
	}

	return helpers.Response(c, 200, message, fmt.Sprintf("Kami telah mengirim email ke %s, silakan cek secara berkala", email))
}

func (ah *AuthController) ForgotPasswordWithAsync(c echo.Context) error {
	var (
		ctx  = c.Request().Context()
		body map[string]interface{}
	)

	if err := json.NewDecoder(c.Request().Body).Decode(&body); err != nil {
		return helpers.Response(c, http.StatusBadRequest, nil, err.Error())
	}

	email := body["email"].(string)
	user, err := ah.userModel.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helpers.Response(c, http.StatusNotFound, nil, "Email tidak ditemukan")
		}
		return helpers.Response(c, http.StatusInternalServerError, nil, err.Error())
	}

	message, err := helpers.EncryptMessageRSA(user.ID.String())
	if err != nil {
		return helpers.Response(c, http.StatusInternalServerError, nil, err.Error())
	}

	// send email using sendgrid
	/*
		go helpers.SendMailSendgrid(message, "Lupa Password", user.Name, user.Email)
	*/

	// send email using gmail
	go helpers.SendMailGmailWithAsync(message, "Lupa Password", user.Email)

	return helpers.Response(c, 200, message, fmt.Sprintf("Kami telah mengirim email ke %s, silakan cek secara berkala", email))
}

func (ah *AuthController) ForgotPasswordWithRabbitMQ(c echo.Context) error {
	var (
		ctx  = c.Request().Context()
		body map[string]interface{}
	)

	if err := json.NewDecoder(c.Request().Body).Decode(&body); err != nil {
		return helpers.Response(c, http.StatusBadRequest, nil, err.Error())
	}

	email := body["email"].(string)
	user, err := ah.userModel.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helpers.Response(c, http.StatusNotFound, nil, "Email tidak ditemukan")
		}
		return helpers.Response(c, http.StatusInternalServerError, nil, err.Error())
	}

	message, err := helpers.EncryptMessageRSA(user.ID.String())
	if err != nil {
		return helpers.Response(c, http.StatusInternalServerError, nil, err.Error())
	}

	payload := structs.Mail{
		TargetName:  user.Name,
		TargetEmail: user.Email,
		Subject:     "Lupa Password",
		Body:        message,
	}

	status, response, data := rabbitmq.RequestCommand(
		rabbitmq.SendMailSendgrid,
		"",
		payload,
		true,
	)

	fmt.Println("status:", status)
	fmt.Println("response:", response)
	fmt.Println("data:", data)

	// send email using gmail
	// err = helpers.SendMailGmail(message, "Lupa Password", user.Email)
	// if err != nil {
	// 	return helpers.Response(c, http.StatusInternalServerError, nil, err.Error())
	// }

	return helpers.Response(c, 200, message, fmt.Sprintf("Kami telah mengirim email ke %s, silakan cek secara berkala", email))
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
