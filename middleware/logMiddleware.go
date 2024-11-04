package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"simple-crud-rnd/helpers"
	"simple-crud-rnd/structs"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type CustomResponseWriter struct {
	echo.Response
	Body *bytes.Buffer
}

func (w *CustomResponseWriter) Write(b []byte) (int, error) {
	w.Body.Write(b) // capture the response body
	return w.Response.Writer.Write(b)
}

var (
	currentLogDate string
	logger         *logrus.Logger
)

func LogMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		payload := getPayload(c)
		checkAndUpdateLogFile()

		res := &CustomResponseWriter{
			Response: *c.Response(),
			Body:     new(bytes.Buffer),
		}
		c.Response().Writer = res
		err := next(c)

		// log the response
		var response interface{}
		if err == nil {
			json.Unmarshal(res.Body.Bytes(), &response)
		}

		user := c.Get("user")
		entry := structs.LogEntry{
			URL:      c.Request().URL.String(),
			Path:     c.Request().Method,
			IP:       c.RealIP(),
			User:     user,
			Body:     payload,
			Response: response,
		}

		logger.WithFields(logrus.Fields{
			"log": entry,
		}).Info("Request log")

		if err := helpers.InsertLogDataToMongoDB(c.Request().Context(), "db_test", "logs", &entry); err != nil {
			log.Printf("Error inserting log to MongoDB: %v", err)
		}

		return err
	}
}

// init new logger
func initLogger() {
	currentLogDate = time.Now().Format("20060102")
	filename := fmt.Sprintf("./logs/app_%s.log", currentLogDate)

	logger = logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{TimestampFormat: time.RFC3339})
	logger.SetOutput(&lumberjack.Logger{
		Filename:   filename,
		MaxSize:    10,   // Maximum size in MB before rotation
		MaxBackups: 7,    // Number of backups to keep
		MaxAge:     30,   // Maximum age of logs in days
		Compress:   true, // Compress old log files
	})
}

// checks if the date has changed to update the log file
func checkAndUpdateLogFile() {
	if currentLogDate == "" {
		initLogger()
		return
	}

	today := time.Now().Format("20060102")
	if today != currentLogDate {
		initLogger()
	}
}
