package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"simple-crud-rnd/config"
	"simple-crud-rnd/helpers"
	"simple-crud-rnd/models"
	rmq "simple-crud-rnd/rabbitmq"
	"simple-crud-rnd/structs"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type UserController struct {
	db          *gorm.DB
	rds         *redis.Client
	model       *models.UserModel
	jobModel    *models.JobModel
	cfg         *config.Config
	imageHelper *helpers.ImageHelper
	assetPath   string
}

func NewUserController(db *gorm.DB, rds *redis.Client, model *models.UserModel, jobModel *models.JobModel, cfg *config.Config, imageHelper *helpers.ImageHelper, assetPath string) *UserController {
	return &UserController{db, rds, model, jobModel, cfg, imageHelper, assetPath}
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
		return helpers.Response(c, http.StatusBadRequest, nil, "Invalid user ID")
	}

	userKey := fmt.Sprintf("user_%v", id)
	cachedData, err := uh.rds.Get(ctx, userKey).Result()
	if err == nil {
		var data structs.User
		if err := json.Unmarshal([]byte(cachedData), &data); err == nil {
			return helpers.Response(c, http.StatusOK, data, "")
		} else {
			log.Printf("Error unmarshaling user data from Redis for user ID %v: %v", id, err)
		}
	} else if err != redis.Nil {
		log.Printf("Error fetching user data from Redis for user ID %v: %v", id, err)
	}

	// If data is not found in Redis, fetch from the database
	data, err := uh.model.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helpers.Response(c, http.StatusNotFound, nil, "Data tidak ditemukan")
		}
		return helpers.Response(c, http.StatusInternalServerError, nil, err.Error())
	}

	// set data user to redis
	rdsErr := helpers.SetRedisJSONCache(ctx, uh.rds, fmt.Sprintf("user_%v", data.ID), data, time.Duration(uh.cfg.Redis.TTL)*time.Second)
	if rdsErr != nil {
		log.Printf("Error setting user data in Redis for user ID %v: %v", data.ID, rdsErr)
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

	rdsErr := helpers.SetRedisJSONCache(ctx, uh.rds, fmt.Sprintf("user_%v", data.ID), data, time.Duration(uh.cfg.Redis.TTL)*time.Second)
	if rdsErr != nil {
		log.Printf("Error setting user data in Redis for user ID %v: %v", data.ID, rdsErr)
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

	rdsErr := helpers.SetRedisJSONCache(ctx, uh.rds, fmt.Sprintf("user_%v", data.ID), data, time.Duration(uh.cfg.Redis.TTL)*time.Second)
	if rdsErr != nil {
		log.Printf("Error setting user data in Redis for user ID %v: %v", data.ID, rdsErr)
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

	rdsErr := uh.rds.Del(ctx, fmt.Sprintf("user_%v", id)).Err()
	if rdsErr != nil {
		log.Printf("Error deleting user data in Redis for user ID %v: %v", id, rdsErr)
	}

	return helpers.Response(c, http.StatusOK, true, "Berhasil hapus user")
}

func (uh *UserController) ExportDataUserToCsv(c echo.Context) error {
	var (
		ctx = c.Request().Context()
		to  = c.QueryParam("to")
	)

	limit := 10
	offset := (1 - 1) * limit
	data, _, err := uh.model.GetAll(ctx, limit, offset)
	if err != nil {
		return helpers.Response(c, http.StatusInternalServerError, data, err.Error())
	}

	filename := fmt.Sprintf("%s/users_%s.csv", "./assets", "users"+time.Now().Format("20061021545"))
	err = helpers.ExportUsersToCSV(filename, data, []string{"Name", "Email", "PhoneNumber", "CreatedAt"})
	if err != nil {
		return helpers.Response(c, http.StatusInternalServerError, data, err.Error())
	}

	job := structs.Job{
		ID:       uuid.NewString(),
		JobName:  "Export Data User",
		Payload:  filename,
		Status:   structs.JOB_PROGRESS,
		Attempts: 0,
	}

	err = uh.jobModel.Create(ctx, &job)
	if err != nil {
		return helpers.Response(c, http.StatusInternalServerError, nil, err.Error())
	}

	// send email using gmail
	go func() {
		err := helpers.SendMailGmail("Berikut adalah file export data users", job.JobName, to, filename)
		if err != nil {
			log.Printf("=> send email to %s is failure. error: %v", to, err.Error())
			updateErr := config.DB.Model(&structs.Job{}).
				Where("id = ?", job.ID).
				Update("status", structs.JOB_FAILED).Error
			if updateErr != nil {
				log.Printf("Failed to update job status to failed for job ID: %s. Error: %v", job.ID, updateErr)
			}
		}

		updateErr := config.DB.Model(&structs.Job{}).
			Where("id = ?", job.ID).
			Update("status", structs.JOB_SUCCESS).Error
		if updateErr != nil {
			log.Printf("Failed to update job status to failed for job ID: %s. Error: %v", job.ID, updateErr)
		}
		log.Printf("send email to %s is successfully", to)
	}()

	return helpers.Response(c, http.StatusOK, job, fmt.Sprintf("Permintaan sedang diproses dan file akan kami kirim ke email %s. Silakan cek secara berkala", to))
}

func (uh *UserController) ExportDataUserToCsvWithRabbitMQ(c echo.Context) error {
	var (
		ctx = c.Request().Context()
		to  = c.QueryParam("to")
	)

	limit := 10
	offset := (1 - 1) * limit
	data, _, err := uh.model.GetAll(ctx, limit, offset)
	if err != nil {
		return helpers.Response(c, http.StatusInternalServerError, data, err.Error())
	}

	filename := fmt.Sprintf("%s/users_%s.csv", "./assets", "users"+time.Now().Format("20061021545"))
	err = helpers.ExportUsersToCSV(filename, data, []string{"Name", "Email", "PhoneNumber", "CreatedAt"})
	if err != nil {
		return helpers.Response(c, http.StatusInternalServerError, data, err.Error())
	}

	job := structs.Job{
		ID:       uuid.NewString(),
		JobName:  "Export Data User",
		Payload:  filename,
		Status:   structs.JOB_PROGRESS,
		Attempts: 0,
	}

	err = uh.jobModel.Create(ctx, &job)
	if err != nil {
		return helpers.Response(c, http.StatusInternalServerError, nil, err.Error())
	}

	jobString, err := json.Marshal(job)
	if err != nil {
		return helpers.Response(c, http.StatusInternalServerError, nil, err.Error())
	}

	// send email using gmail
	if err := rmq.SendMessage(uh.cfg, &rmq.PublishMessage{
		Key:     helpers.KEY_EXPORT_DATA_USER,
		Payload: to + "|||" + string(jobString),
	}); err != nil {
		return helpers.Response(c, http.StatusInternalServerError, nil, err.Error())
	}

	return helpers.Response(c, http.StatusOK, job, fmt.Sprintf("Permintaan sedang diproses dan file akan kami kirim ke email %s. Silakan cek secara berkala", to))
}
