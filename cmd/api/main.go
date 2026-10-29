package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/magomed066/go-task-queue/internal/env"
)

func initialize() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	initialize()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := config{
		db: dbConfig{
			dns: env.GetString("DB_URL", ""),
		},
		logLevel:  env.GetString("LOG_LEVEL", "info"),
		logFormat: env.GetString("LOG_FORMAT", "json"),

		taskTimeout:   env.GetTime("TASK_TIMEOUT", time.Second*30),
		pollInterval:  env.GetTime("POLL_INTERVAL", time.Second*1),
		leaseDuration: env.GetTime("LEASE_DURATION", time.Second*45),

		workerPoolSize: env.GetInt("WORKER_POOL_SIZE", 10),
		maxRetries:     env.GetInt("MAX_RETRIES", 5),
	}

	app := application{
		addr:   env.GetString("ADDR", ":8080"),
		config: cfg,
	}

	err := app.run(ctx, app.mount())
	if err != nil {
		fmt.Println("Exit error:", err)
		os.Exit(1)
	}

}
