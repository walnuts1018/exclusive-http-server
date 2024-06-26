package main

import (
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lmittmann/tint"
)

func main() {
	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{
		TimeFormat: time.RFC3339,
	}))
	slog.SetDefault(logger)

	sleepDuration := 0 * time.Second
	sleepDurationEnv := os.Getenv("SLEEP_DURATION")
	if sleepDurationEnv != "" {
		parsed, err := time.ParseDuration(sleepDurationEnv)
		if err != nil {
			slog.Info(fmt.Sprintf("Failed to parse SLEEP_DURATION, using default: %v", err))
		} else {
			sleepDuration = parsed
		}
	}
	slog.Info(fmt.Sprintf("Using sleep time: %v", sleepDuration))

	counter := 0
	mutex := sync.Mutex{}

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		mutex.Lock()
		defer mutex.Unlock()
		counter++
		time.Sleep(sleepDuration)
		c.JSON(200, gin.H{
			"counter": counter,
		})
		slog.Info(fmt.Sprintf("Counter: %d", counter))
	})
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	slog.Info("Starting server on port 8080")
	if err := r.Run(":8080"); err != nil {
		slog.Error(fmt.Sprintf("Failed to start server: %v", err))
	}
}
