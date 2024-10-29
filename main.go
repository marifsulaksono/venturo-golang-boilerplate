package main

import (
	"fmt"
	"log"

	"simple-crud-rnd/config"
	"simple-crud-rnd/helpers"
	"simple-crud-rnd/routes"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Error loading configs")
	}

	db, err := config.InitDatabase(cfg)
	if err != nil {
		log.Fatalln("Error opening database")
	}

	id := "c1daac6e-47fc-4200-a0e9-58832e8c5de0"
	str, err := helpers.EncryptMessageRSA(id)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(str)

	e := routes.NewHTTPServer(cfg, db)
	e.RunHTTPServer()
}
