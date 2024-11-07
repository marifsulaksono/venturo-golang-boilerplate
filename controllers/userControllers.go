package controllers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"simple-crud-rnd/config"
	"simple-crud-rnd/helpers"
	"simple-crud-rnd/models"
	"simple-crud-rnd/structs"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type UserController struct {
	db          *gorm.DB
	model       *models.UserModel
	cfg         *config.Config
	imageHelper *helpers.ImageHelper
	assetPath   string
}

func NewUserController(db *gorm.DB, model *models.UserModel, cfg *config.Config, imageHelper *helpers.ImageHelper, assetPath string) *UserController {
	return &UserController{db, model, cfg, imageHelper, assetPath}
}

func (uh *UserController) Index(c echo.Context) error {
	var ctx = c.Request().Context()
	per_page, err := strconv.Atoi(c.QueryParam("per_page"))
	if err != nil {
		per_page = 10
		log.Printf("Failed to parse per_page query parameter. Defaulting to %d", per_page)
	}
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		page = 1
		log.Printf("Failed to parse page query parameter. Defaulting to %d", page)
	}

	offset := (page - 1) * per_page
	data, total, err := uh.model.GetAll(ctx, per_page, offset)
	if err != nil {
		return helpers.Response(c, http.StatusInternalServerError, data, err.Error())
	}
	pagedData := helpers.PageData(data, total)
	return helpers.Response(c, http.StatusOK, pagedData, "")
}

func (uh *UserController) GetById(c echo.Context) error {
	var ctx = c.Request().Context()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return helpers.SendTraceErrorToSentry(err)
	}

	data, err := uh.model.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helpers.Response(c, http.StatusNotFound, nil, "Data tidak ditemukan")
		}
		return helpers.Response(c, http.StatusInternalServerError, nil, err.Error())
	}
	return helpers.Response(c, http.StatusOK, data, "")
}

func (uh *UserController) Create(c echo.Context) error {
	var (
		ctx     = c.Request().Context()
		request structs.UserRequest
	)

	if err := c.Bind(&request); err != nil {
		return helpers.Response(c, http.StatusBadRequest, nil, err.Error())
	}

	if err := c.Validate(request); err != nil {
		return helpers.Response(c, http.StatusBadRequest, nil, err.Error())
	}

	if request.Photo != "" {
		photo_url, err := uh.imageHelper.Writer(request.Photo, fmt.Sprintf("%s.png", time.Now().Format("20061021545.000000000")))
		if err != nil {
			return helpers.Response(c, http.StatusInternalServerError, nil, err.Error())
		}
		request.Photo = photo_url
	}

	data, err := uh.model.Create(ctx, &request)
	if err != nil {
		return helpers.Response(c, http.StatusInternalServerError, nil, err.Error())
	}

	return helpers.Response(c, http.StatusCreated, data, "Berhasil simpan user")
}

func (uh *UserController) Update(c echo.Context) error {
	var (
		ctx     = c.Request().Context()
		request structs.User
	)

	if err := c.Bind(&request); err != nil {
		return helpers.Response(c, http.StatusBadRequest, nil, err.Error())
	}

	if err := c.Validate(request); err != nil {
		return helpers.Response(c, http.StatusBadRequest, nil, err.Error())
	}

	user, err := uh.model.GetById(ctx, request.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helpers.Response(c, http.StatusNotFound, nil, "Data tidak ditemukan")
		}
		return helpers.Response(c, http.StatusInternalServerError, nil, err.Error())
	}

	if request.Photo != "" {
		photo_url, err := uh.imageHelper.Writer(request.Photo, fmt.Sprintf("%s.png", time.Now().Format("20061021545.000000000")))
		if err != nil {
			return helpers.Response(c, http.StatusInternalServerError, nil, err.Error())
		}
		request.Photo = photo_url
	}

	data, err := uh.model.Update(ctx, &request)
	if err != nil {
		return helpers.Response(c, http.StatusInternalServerError, nil, err.Error())
	}

	if user.Photo != "" {
		err = helpers.MoveToTrash(user.Photo)
		if err != nil {
			log.Printf("Failed to move file %s to trash: %s", user.Photo, err)
		}
	}

	return helpers.Response(c, http.StatusOK, data, "Berhasil update user")
}

func (uh *UserController) Delete(c echo.Context) error {
	var ctx = c.Request().Context()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return err
	}
	if err := uh.model.Delete(ctx, id); err != nil {
		if errors.Is(err, errors.New("no rows deleted")) {
			return helpers.Response(c, http.StatusNotFound, nil, "Data tidak ditemukan, gagal menghapus")
		}
		return helpers.Response(c, http.StatusInternalServerError, nil, err.Error())
	}

	return helpers.Response(c, http.StatusOK, true, "Berhasil hapus user")
}
