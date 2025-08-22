package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	dataaccess "github.com/AmirMahdyJebreily/timeline-example/internal/data"
	"github.com/AmirMahdyJebreily/timeline-example/internal/timeline"
)

const (
	DEFAULT_LIMIT  uint = 10
	DEFAULT_OFFSET uint = 0
)

func CreatePost(tl *timeline.TimelineService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var post dataaccess.Post

		if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
			http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
			return
		}

		ctx := r.Context()

		postID, err := tl.NewPost(ctx, post)
		if err != nil {
			http.Error(w, fmt.Sprintf("Insert failed: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		if err := json.NewEncoder(w).Encode(postID); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}

func GetTimeline(tl *timeline.TimelineService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		query := r.URL.Query()

		limitStr := query.Get("limit")
		limit := DEFAULT_LIMIT
		if limitStr != "" {
			l, err := strconv.Atoi(limitStr)
			if err != nil || l < 1 || l > 100 {
				http.Error(w, "Invalid 'limit' parameter", http.StatusBadRequest)
				return
			}
			limit = uint(l)
		}

		offsetStr := query.Get("offset")
		offset := DEFAULT_OFFSET
		if offsetStr != "" {
			o, err := strconv.Atoi(offsetStr)
			if err != nil || o < 0 {
				http.Error(w, "Invalid 'offset' parameter", http.StatusBadRequest)
				return
			}
			offset = uint(o)
		}

		userIdStr := query.Get("userId")
		if userIdStr == "" {
			http.Error(w, "Missing userId parameter", http.StatusBadRequest)
			return
		}
		userId, err := strconv.Atoi(userIdStr)
		if err != nil || userId < 1 {
			http.Error(w, "Invalid userId parameter", http.StatusBadRequest)
			return
		}

		ctx := r.Context()

		posts, err := tl.GetTimeLine(ctx, uint(userId), limit, offset)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(posts); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}
