package controllers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"simple-crud-rnd/config"
	"simple-crud-rnd/helpers"
	"simple-crud-rnd/models"
	"simple-crud-rnd/structs"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type RoleController struct {
	db    *gorm.DB
	model *models.RoleModel
	cfg   *config.Config
}

func NewRoleController(db *gorm.DB, model *models.RoleModel, cfg *config.Config) *RoleController {
	return &RoleController{db, model, cfg}
}

func (uh *RoleController) Index(c echo.Context) error {
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

func (uh *RoleController) GetById(c echo.Context) error {
	var ctx = c.Request().Context()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return err
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

func (uh *RoleController) Create(c echo.Context) error {
	var (
		ctx     = c.Request().Context()
		request structs.RoleRequest
	)

	if err := c.Bind(&request); err != nil {
		return helpers.Response(c, http.StatusBadRequest, nil, err.Error())
	}

	if err := c.Validate(request); err != nil {
		return helpers.Response(c, http.StatusBadRequest, nil, err.Error())
	}

	access, err := json.Marshal(request.Access)
	if err != nil {
		return helpers.Response(c, http.StatusInternalServerError, nil, err.Error())
	}

	payload := structs.Role{
		Name:   request.Name,
		Access: string(access),
	}
	request.Access = string(access)

	data, err := uh.model.Create(ctx, &payload)
	if err != nil {
		return helpers.Response(c, http.StatusInternalServerError, nil, err.Error())
	}

	return helpers.Response(c, http.StatusCreated, data, "Berhasil simpan role")
}

func (uh *RoleController) Update(c echo.Context) error {
	var (
		ctx     = c.Request().Context()
		request structs.RoleRequest
	)

	if err := c.Bind(&request); err != nil {
		return helpers.Response(c, http.StatusBadRequest, nil, err.Error())
	}

	if err := c.Validate(request); err != nil {
		return helpers.Response(c, http.StatusBadRequest, nil, err.Error())
	}

	access, err := json.Marshal(request.Access)
	if err != nil {
		return helpers.Response(c, http.StatusInternalServerError, nil, err.Error())
	}

	payload := structs.Role{
		ID:     request.ID,
		Name:   request.Name,
		Access: string(access),
	}

	data, err := uh.model.Update(ctx, &payload)
	if err != nil {
		return helpers.Response(c, http.StatusInternalServerError, nil, err.Error())
	}

	return helpers.Response(c, http.StatusOK, data, "Berhasil update role")
}

func (uh *RoleController) Delete(c echo.Context) error {
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

	return helpers.Response(c, http.StatusOK, true, "Berhasil hapus role")
}
