package controllers

import (
	"errors"
	"log"
	"net/http"
	"simple-crud-rnd/config"
	"simple-crud-rnd/helpers"
	"simple-crud-rnd/models"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type JobController struct {
	db       *gorm.DB
	jobModel *models.JobModel
	cfg      *config.Config
}

func NewJobController(db *gorm.DB, jobModel *models.JobModel, cfg *config.Config) *JobController {
	return &JobController{db, jobModel, cfg}
}

func (jh *JobController) Index(c echo.Context) error {
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
	data, total, err := jh.jobModel.GetAll(ctx, per_page, offset)
	if err != nil {
		return helpers.Response(c, http.StatusInternalServerError, data, err.Error())
	}
	pagedData := helpers.PageData(data, total)
	return helpers.Response(c, http.StatusOK, pagedData, "")
}

func (jh *JobController) GetById(c echo.Context) error {
	var (
		ctx = c.Request().Context()
		id  = c.Param("id")
	)

	data, err := jh.jobModel.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helpers.Response(c, http.StatusNotFound, nil, "Data tidak ditemukan")
		}
		return helpers.Response(c, http.StatusInternalServerError, nil, err.Error())
	}
	return helpers.Response(c, http.StatusOK, data, "")
}
