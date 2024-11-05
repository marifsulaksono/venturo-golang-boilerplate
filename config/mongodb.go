package config

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoCLI *mongo.Client

// using mongodb atlas https://cloud.mongodb.com/
func InitMongoDB(cfg *Config) (*mongo.Client, error) {
	url := fmt.Sprintf("mongodb+srv://%s:%s@%s.mongodb.net/?retryWrites=true&w=majority", cfg.MongoDB.Username, cfg.MongoDB.Password, cfg.MongoDB.Cluster)
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(url).SetServerAPIOptions(serverAPI)

	// connect to MongoDB
	db, err := mongo.Connect(context.TODO(), opts.SetTLSConfig(&tls.Config{InsecureSkipVerify: true}))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// test mongodb connection
	if err := db.Ping(context.TODO(), nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	MongoCLI = db

	log.Println("Succees to connect to mongodb")
	return db, nil
}

// using mongodb shell (local)
func InitLocalMongoDB(cfg *Config) (*mongo.Client, error) {
	url := fmt.Sprintf("mongodb://%s:%s", cfg.MongoDB.Host, cfg.MongoDB.Port)
	opts := options.Client().ApplyURI(url)
	opts.SetAuth(options.Credential{
		AuthSource: cfg.MongoDB.Name,
		Username:   cfg.MongoDB.Username,
		Password:   cfg.MongoDB.Password,
	})

	// connect to MongoDB
	db, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// test MongoDB connection
	if err := db.Ping(context.TODO(), nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	MongoCLI = db

	log.Println("Successfully connected to MongoDB")
	return db, nil
}
