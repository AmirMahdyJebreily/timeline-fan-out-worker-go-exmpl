package handlers

import (
	"fmt"
	"net/http"
)

func CreatePost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			q := r.URL.Query()
			senderId := q.Get("sender_id")
			content := q.Get("content")

			// TODO:Connect to sql

			r.Response.StatusCode = http.StatusCreated
			fmt.Fprintf(w, "Created!")
			r.Response.Write(w)
		}
	}
}

func GetTimeline() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Not implemented yet!")
	}
}
