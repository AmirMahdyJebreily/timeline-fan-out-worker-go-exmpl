package dataaccess

import "time"

type Post struct {
	ID        uint      `db:"id" json:"id,omitempty"`
	SenderID  uint      `db:"sender_id" json:"sender_id"`
	Content   string    `db:"content" json:"content"`
	CreatedAt time.Time `db:"created_at" json:"created_at,omitempty"`
}

type SubscriberUser struct {
	SenderID     uint `db:"sender_id"`
	SubscriberID uint `db:"subscriber_id"`
}
