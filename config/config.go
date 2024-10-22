package config

import (
	"log"
	"os"
	"simple-crud-rnd/structs"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type (
	Config struct {
		Database     Database
		HTTP         HTTP
		JWT          JWT
		AssetStorage AssetStorage
		RabbitMQ     RabbitMQConf
	}
	Database struct {
		Username string
		Password string
		Host     string
		Port     string
		Name     string
	}
	HTTP struct {
		Host          string
		Port          int
		Domain        string
		AssetEndpoint string
	}
	JWT struct {
		Secret []byte
		Config echojwt.Config
	}
	AssetStorage struct {
		Path string
	}
	RabbitMQConf struct {
		Username string
		Password string
		Host     string
		Port     string
		Vhost    string
		Queue    string
	}
)

func LoadConfig() (*Config, error) {
	errEnv := godotenv.Load()

	if errEnv != nil {
		log.Fatal("Unable to load .env file")
	}

	dbUsername, _ := configDefaults("DB_USERNAME", "mysql")
	dbPassword, _ := configDefaults("DB_PASSWORD", "changeme")
	dbHost, _ := configDefaults("DB_HOST", "127.0.0.1")
	dbPort, _ := configDefaults("DB_PORT", "3306")
	dbName, _ := configDefaults("DB_NAME", "mysql")

	listenHost, _ := configDefaults("LISTEN_HOST", "127.0.0.1")
	listenPort, _ := configDefaults("LISTEN_PORT", "8080")
	intListenPort, err := strconv.Atoi(listenPort)
	if err != nil {
		log.Fatal("Port must be a number")
	}
	domain, _ := configDefaults("DOMAIN", "http://localhost")
	assetPath, _ := configDefaults("ASSET_PATH", "api/v1/assets")
	jwtSecret, _ := configDefaults("JWT_SECRET", "")
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(structs.JWTUser)
		},
		SigningKey: []byte(jwtSecret),
	}
	storagePath, _ := configDefaults("ASSET_PATH", "./")

	rmqUser, _ := configDefaults("RQ_USERNAME", "guest")
	rmqPass, _ := configDefaults("RQ_PASSWORD", "guest")
	rmqHost, _ := configDefaults("RQ_HOST", "127.0.0.1")
	rmqPort, _ := configDefaults("RQ_PORT", "5672")
	rmqVhost, _ := configDefaults("RQ_VHOST", "/")
	rmqQueue, _ := configDefaults("RQ_QUEUE", "my_queue")

	var cfg Config = Config{
		Database: Database{
			Username: dbUsername,
			Password: dbPassword,
			Host:     dbHost,
			Port:     dbPort,
			Name:     dbName,
		},
		HTTP: HTTP{
			Host:          listenHost,
			Port:          intListenPort,
			Domain:        domain,
			AssetEndpoint: assetPath,
		},
		JWT: JWT{
			Secret: []byte(jwtSecret),
			Config: config,
		},
		AssetStorage: AssetStorage{
			Path: storagePath,
		},
		RabbitMQ: RabbitMQConf{
			Username: rmqUser,
			Password: rmqPass,
			Host:     rmqHost,
			Port:     rmqPort,
			Vhost:    rmqVhost,
			Queue:    rmqQueue,
		},
	}

	return &cfg, nil
}

func configDefaults(env, defaults string) (string, bool) {
	value, ok := os.LookupEnv(env)
	if !ok {
		log.Printf("%s is unset. Resorting to default value of %s", env, defaults)
		return defaults, ok
	}
	return value, ok
}
