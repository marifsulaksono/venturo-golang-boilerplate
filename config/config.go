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
		Sentry       Sentry
		HTTP         HTTP
		JWT          JWT
		AssetStorage AssetStorage
		MongoDB      MongoDB
	}
	Database struct {
		Username string
		Password string
		Host     string
		Port     string
		Name     string
	}
	Sentry struct {
		Dsn string
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
	MongoDB struct {
		Host     string
		Port     string
		Username string
		Password string
		Cluster  string
		Name     string
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

	sentryDsn, _ := configDefaults("SENTRY_DSN", "")

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

	mongoHost, _ := configDefaults("MONGODB_HOST", "127.0.0.1")
	mongoPort, _ := configDefaults("MONGODB_PORT", "27017")
	mongoUser, _ := configDefaults("MONGODB_USER", "user")
	mongoPass, _ := configDefaults("MONGODB_PASS", "user")
	mongoCluster, _ := configDefaults("MONGODB_CLUSTER", "cluster0")
	mongoName, _ := configDefaults("MONGODB_NAME", "mongo_database")

	var cfg Config = Config{
		Database: Database{
			Username: dbUsername,
			Password: dbPassword,
			Host:     dbHost,
			Port:     dbPort,
			Name:     dbName,
		},
		Sentry: Sentry{
			Dsn: sentryDsn,
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
		MongoDB: MongoDB{
			Host:     mongoHost,
			Port:     mongoPort,
			Username: mongoUser,
			Password: mongoPass,
			Cluster:  mongoCluster,
			Name:     mongoName,
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
