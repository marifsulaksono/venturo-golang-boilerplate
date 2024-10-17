package config

import (
	"fmt"
	"log"
	"simple-crud-rnd/structs"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDatabase(cfg *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Succees to connect to database")

	if err := db.AutoMigrate(
		&structs.User{},
	); err != nil {
		log.Fatal("Failed to migrate to database:", err)
	}

	log.Println("Succees to migrate to database")

	return db, err
}
