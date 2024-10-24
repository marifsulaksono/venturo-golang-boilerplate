package routes

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"simple-crud-rnd/config"
	"simple-crud-rnd/helpers"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type HTTPServer struct {
	db         *gorm.DB
	cfg        *config.Config
	httpServer *echo.Echo
}

func NewHTTPServer(cfg *config.Config, db *gorm.DB) HTTPServer {
	e := echo.New()
	e.Validator = helpers.NewValidator(validator.New())

	return HTTPServer{
		db:         db,
		cfg:        cfg,
		httpServer: e,
	}
}

func testPort(port int) (int, error) {
	ln, err := net.Listen("tcp", ":"+fmt.Sprint(port))
	if err == nil {
		_ = ln.Close()
		return port, nil
	}
	log.Println("port", port, "is in use. Searching for next free port.")

	for i := port; i <= 65535; i++ {
		nextLn, err := net.Listen("tcp", ":"+fmt.Sprint(i))
		if err == nil {
			_ = nextLn.Close()
			return i, nil
		}
	}
	return 0, errors.New("No free available ports")
}

func (s *HTTPServer) RunHTTPServer() {
	api := InitVersionOne(s.httpServer, s.db, s.cfg)

	s.httpServer.Static(api.cfg.HTTP.AssetEndpoint, api.cfg.AssetStorage.Path)
	api.UserAndAuth()
	api.Role()
	api.Job()

	openPort, err := testPort(s.cfg.HTTP.Port)
	if err != nil {
		log.Fatal(err)
	}

	if err := s.httpServer.Start(fmt.Sprintf(":%d", openPort)); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
