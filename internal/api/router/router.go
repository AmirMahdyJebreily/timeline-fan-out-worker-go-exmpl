package router

import (
	"log"
	"net/http"
	"time"

	"github.com/AmirMahdyJebreily/timeline-example/internal/api/handlers"
	"github.com/AmirMahdyJebreily/timeline-example/internal/database"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Incoming request: %s %s", r.Method, r.URL.Path)

		next.ServeHTTP(w, r)

		log.Printf("Request %s %s processed in %v", r.Method, r.URL.Path, time.Since(start))
	})
}

func New(repo database.PostsRepository) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/posts", handlers.CreatePost(repo))
	mux.HandleFunc("/timeline", handlers.GetTimeline())

	var handler http.Handler = mux

	handler = loggingMiddleware(handler)

	return handler
}
