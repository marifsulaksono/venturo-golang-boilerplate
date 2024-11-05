package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"simple-crud-rnd/config"
	"simple-crud-rnd/routes"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Error loading configs:", err)
	}

	db, err := config.InitDatabase(cfg)
	if err != nil {
		log.Fatalln("Error opening database:", err)
	}

	if err := config.InitSentry(cfg); err != nil {
		log.Fatalln("Error initialize sentry:", err)
	}

	_, err = config.InitLocalMongoDB(cfg)
	if err != nil {
		log.Fatalln("Error connect mongodb:", err)
	}

	e := routes.NewHTTPServer(cfg, db)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		e.RunHTTPServer()
	}()

	<-sig

	log.Println("Shutting down.....")
	config.FlushSentry()
	log.Println("Gracefully terminated.")
}
