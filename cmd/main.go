package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/AmirMahdyJebreily/timeline-example/internal/api/router"
	dataAccess "github.com/AmirMahdyJebreily/timeline-example/internal/data"
)

func main() {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = ":8080" // fallback default
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

	dataAccesser := dataAccess.New(db)
	r := router.New(dataAccesser)

	server := &http.Server{
		Addr:    port,
		Handler: r,
	}

	fmt.Printf("Server is running on port %v...\n", port)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
