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

	sleepTime := 100 * time.Millisecond
	sleeptimeEnv := os.Getenv("SLEEP_TIME")
	if sleeptimeEnv != "" {
		parsed, err := time.ParseDuration(sleeptimeEnv)
		if err != nil {
			slog.Info(fmt.Sprintf("Failed to parse SLEEP_TIME, using default: %v", err))
		} else {
			sleepTime = parsed
		}
	}
	slog.Info(fmt.Sprintf("Using sleep time: %v", sleepTime))

	counter := 0
	mutex := sync.Mutex{}

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		mutex.Lock()
		defer mutex.Unlock()
		counter++
		time.Sleep(sleepTime)
		c.JSON(200, gin.H{
			"counter": counter,
		})
	})

	slog.Info("Starting server on port 8080")
	if err := r.Run(":8080"); err != nil {
		slog.Error(fmt.Sprintf("Failed to start server: %v", err))
	}
}
