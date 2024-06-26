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

	counter := 0
	mutex := sync.Mutex{}

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		mutex.Lock()
		defer mutex.Unlock()
		counter++
		time.Sleep(100 * time.Millisecond)
		c.JSON(200, gin.H{
			"counter": counter,
		})
	})

	slog.Info("Starting server on port 8080")
	if err := r.Run(":8080"); err != nil {
		slog.Error(fmt.Sprintf("Failed to start server: %v", err))
	}
}
