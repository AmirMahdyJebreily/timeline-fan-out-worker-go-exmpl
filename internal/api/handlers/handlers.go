package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AmirMahdyJebreily/timeline-example/internal/database"
	"github.com/AmirMahdyJebreily/timeline-example/internal/database/entities"
)

func CreatePost(repo database.PostsRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var posts []entities.Post

		if err := json.NewDecoder(r.Body).Decode(&posts); err != nil {
			http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
			return
		}

		ctx := r.Context()

		if _, err := repo.InsertPost(ctx, posts); err != nil {
			http.Error(w, fmt.Sprintf("Insert failed: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Created!"))
	}
}

func GetTimeline() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Not implemented yet!")
	}
}
