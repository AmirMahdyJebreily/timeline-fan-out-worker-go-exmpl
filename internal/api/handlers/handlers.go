package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	dataaccess "github.com/AmirMahdyJebreily/timeline-example/internal/data"
)

func CreatePost(da dataaccess.DataAccesser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var posts []dataaccess.Post

		if err := json.NewDecoder(r.Body).Decode(&posts); err != nil {
			http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
			return
		}

		ctx := r.Context()

		postIDs, err := da.BulkInsertPosts(ctx, posts)
		if err != nil {
			http.Error(w, fmt.Sprintf("Insert failed: %v", err), http.StatusInternalServerError)
			return
		}

		response := struct {
			PostIDs []uint `json:"post_ids"`
		}{
			PostIDs: postIDs,
		}

		responseJSON, err := json.Marshal(response)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to marshal response: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(responseJSON)
	}
}

func GetTimeline() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Not implemented yet!")
	}
}
