package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"

	"github.com/AmirMahdyJebreily/timeline-example/internal/api/router"
	"github.com/AmirMahdyJebreily/timeline-example/internal/cache"
	dataAccess "github.com/AmirMahdyJebreily/timeline-example/internal/data"
	"github.com/AmirMahdyJebreily/timeline-example/internal/timeline"
	"github.com/AmirMahdyJebreily/timeline-example/internal/workers"
)

const DEFAULT_APP_PORT = ":8080"

func main() {

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = DEFAULT_APP_PORT
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true&multiStatements=true", dbUser, dbPass, dbHost, dbName)
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}

	dataAccess := dataAccess.New(db)

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
		DB:   0,
	})
	timelineCache := cache.New(redisClient)

	workerPool, err := workers.New(context.Background(), 50, 25, 200)
	if err != nil {
		log.Fatalf("Failed to create singleton worker pool: %v", err)
	}

	if err = workerPool.InitWorkers(false); err != nil {
		log.Fatalf("Failed to initialize the worker pool: %v", err)
	}

	tl := timeline.New(dataAccess, timelineCache, workerPool)

	r := router.New(tl)
	server := &http.Server{
		Addr:    port,
		Handler: r,
	}

	fmt.Printf("Server is running on port %v...\n", port)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
