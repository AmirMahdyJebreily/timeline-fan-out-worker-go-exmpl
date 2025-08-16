package entities

import "time"

type Post struct {
	Id          int
	SenderId    int
	Content     string
	CreatedDate *time.Time
}
