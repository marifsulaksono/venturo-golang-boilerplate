package controllers

import (
	"errors"
	"net/http"
	"simple-crud-rnd/config"
	"simple-crud-rnd/helpers"
	"simple-crud-rnd/models"
	"simple-crud-rnd/structs"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthController struct {
	db    *gorm.DB
	model *models.UserModel
	cfg   *config.Config
}

func NewAuthController(db *gorm.DB, model *models.UserModel, cfg *config.Config) *AuthController {
	return &AuthController{db, model, cfg}
}

func (uh *AuthController) Login(c echo.Context) error {
	var (
		ctx     = c.Request().Context()
		request structs.Login
	)

	if err := c.Bind(&request); err != nil {
		return helpers.Response(c, http.StatusBadRequest, nil, err.Error())
	}

	if err := c.Validate(&request); err != nil {
		return helpers.Response(c, http.StatusBadRequest, nil, err.Error())
	}

	user, err := uh.model.GetByEmail(ctx, request.Email)
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

	accessToken, expiredAt, err := helpers.GenerateTokenJWT(user)
	if err != nil {
		return helpers.Response(c, http.StatusInternalServerError, nil, err.Error())
	}

	response := structs.LoginResponse{
		AccessToken: accessToken,
		ExpiresAt:   *expiredAt,
		Metadata: structs.Metadata{
			Name:   user.Name,
			Email:  user.Email,
			Access: user.Role.Access,
		},
	}

	return helpers.Response(c, http.StatusOK, response, "")
}
