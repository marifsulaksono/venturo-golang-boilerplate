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
	userModel := models.NewUserModel(av.db)
	imageHelper, err := helpers.NewImageHelper(av.cfg.AssetStorage.Path, "profile_photos")
	if err != nil {
		log.Fatal("Failed to initiate an image helper:", err)
	}

	userController := controllers.NewUserController(av.db, userModel, av.cfg, imageHelper, av.assetsPath)
	authController := controllers.NewAuthController(av.db, userModel, av.cfg)

	auth := av.api.Group("/auth")
	auth.POST("/login", authController.Login)
	// auth.POST("/signup", userController.Create)

	// user := av.api.Group("/users", echojwt.WithConfig(av.cfg.JWT.Config))
	user := av.api.Group("/users")

	user.GET("", userController.Index, middleware.RoleMiddleware("roles.view"))
	user.POST("", userController.Create, middleware.RoleMiddleware("roles.create"))
	user.GET("/:id", userController.GetById, middleware.RoleMiddleware("roles.view"))
	user.PUT("", userController.Update, middleware.RoleMiddleware("roles.update"))
	user.DELETE("/:id", userController.Delete, middleware.RoleMiddleware("roles.delete"))
}

func (av *APIVersionOne) Role() {
	roleModel := models.NewRoleModel(av.db)
	roleController := controllers.NewRoleController(av.db, roleModel, av.cfg)

	role := av.api.Group("/roles")
	// role := av.api.Group("/roles", echojwt.WithConfig(av.cfg.JWT.Config))
	role.GET("", roleController.Index, middleware.RoleMiddleware("user.view"))
	role.POST("", roleController.Create)
	role.GET("/:id", roleController.GetById, middleware.RoleMiddleware("user.view"))
	role.PUT("", roleController.Update, middleware.RoleMiddleware("user.update"))
	role.DELETE("/:id", roleController.Delete, middleware.RoleMiddleware("user.delete"))
}
