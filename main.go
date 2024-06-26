package main

import (
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lmittmann/tint"
	"github.com/penglongli/gin-metrics/ginmetrics"
	"golang.org/x/sync/errgroup"
)

func main() {
	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{
		TimeFormat: time.RFC3339,
	}))
	slog.SetDefault(logger)

	podName := os.Getenv("POD_NAME")
	if podName == "" {
		slog.Error("POD_NAME is not set")
	}

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
	appRouter := gin.Default()

	metricRouter := gin.Default()
	m := ginmetrics.GetMonitor()
	m.UseWithoutExposingEndpoint(appRouter)
	m.SetMetricPath("/metrics")
	m.Expose(metricRouter)

	appRouter.GET("/", func(c *gin.Context) {
		mutex.Lock()
		defer mutex.Unlock()
		counter++
		time.Sleep(sleepDuration)
		c.JSON(200, gin.H{
			"counter": counter,
			"pod":     podName,
		})
		slog.Info(fmt.Sprintf("Counter: %d", counter))
	})
	metricRouter.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	eg := errgroup.Group{}
	eg.Go(func() error {
		if err := metricRouter.Run(":8081"); err != nil {
			return fmt.Errorf("failed to start metrics server: %w", err)
		}
		return nil
	})
	eg.Go(func() error {
		if err := appRouter.Run(":8080"); err != nil {
			return fmt.Errorf("failed to start server: %w", err)
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		slog.Error(err.Error())
	}

}
