package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/AmirMahdyJebreily/timeline-example/internal/api/router"
)

func main() {
	port := os.Getenv("APP_PORT")
	router := router.New()

	server := &http.Server{
		Addr:    port,
		Handler: *router,
	}

	fmt.Printf("Server is running on port %v...\n", port)

	// 3. Start the server.
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
