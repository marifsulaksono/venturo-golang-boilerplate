package routes

import (
	"fmt"
	"log"
	"simple-crud-rnd/config"
	"simple-crud-rnd/controllers"
	"simple-crud-rnd/helpers"
	"simple-crud-rnd/middleware"
	"simple-crud-rnd/models"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type APIVersionOne struct {
	e          *echo.Echo
	db         *gorm.DB
	cfg        *config.Config
	api        *echo.Group
	assetsPath string
}

func InitVersionOne(e *echo.Echo, db *gorm.DB, cfg *config.Config) *APIVersionOne {
	return &APIVersionOne{
		e,
		db,
		cfg,
		e.Group("/api/v1"),
		fmt.Sprintf("%s/%s", cfg.HTTP.Domain, cfg.HTTP.AssetEndpoint),
	}
}

func (av *APIVersionOne) UserAndAuth() {
	authModel := models.NewAuthModel(av.db)
	userModel := models.NewUserModel(av.db)
	jobModel := models.NewJobModel(av.db)
	imageHelper, err := helpers.NewImageHelper(av.cfg.AssetStorage.Path, "profile_photos")
	if err != nil {
		log.Fatal("Failed to initiate an image helper:", err)
	}

	authController := controllers.NewAuthController(av.db, authModel, userModel, jobModel, av.cfg)
	userController := controllers.NewUserController(av.db, userModel, av.cfg, imageHelper, av.assetsPath)

	auth := av.api.Group("/auth")
	auth.POST("/login", authController.Login)
	auth.POST("/refresh", authController.RefreshAccessToken)
	auth.POST("/forgot-password", authController.ForgotPasswordWithAsync)
	auth.POST("/logout", authController.Logout)

	user := av.api.Group("/users")
	user.GET("", userController.Index, middleware.RoleMiddleware("user.view"))
	user.POST("", userController.Create)
	user.GET("/:id", userController.GetById, middleware.RoleMiddleware("user.view"))
	user.PUT("", userController.Update, middleware.RoleMiddleware("user.update"))
	user.DELETE("/:id", userController.Delete, middleware.RoleMiddleware("user.delete"))
}

func (av *APIVersionOne) Role() {
	roleModel := models.NewRoleModel(av.db)
	roleController := controllers.NewRoleController(av.db, roleModel, av.cfg)

	role := av.api.Group("/roles")
	role.GET("", roleController.Index, middleware.RoleMiddleware("roles.view"))
	role.POST("", roleController.Create, middleware.RoleMiddleware("roles.create"))
	role.GET("/:id", roleController.GetById, middleware.RoleMiddleware("roles.view"))
	role.PUT("", roleController.Update, middleware.RoleMiddleware("roles.update"))
	role.DELETE("/:id", roleController.Delete, middleware.RoleMiddleware("roles.delete"))
}

// func (av *APIVersionOne) Job() {
// 	jobModel := models.NewJobModel(av.db)
// 	jobController := controllers.NewJobController(av.db, jobModel, av.cfg)

// 	job := av.api.Group("/jobs")
// 	job.GET("", jobController.Index, middleware.RoleMiddleware("jobs.view"))
// 	job.POST("", jobController.Create, middleware.RoleMiddleware("jobs.create"))
// 	job.GET("/:id", jobController.GetById, middleware.RoleMiddleware("jobs.view"))
// 	job.PUT("", jobController.Update, middleware.RoleMiddleware("jobs.update"))
// 	job.DELETE("/:id", jobController.Delete, middleware.RoleMiddleware("jobs.delete"))
// }
