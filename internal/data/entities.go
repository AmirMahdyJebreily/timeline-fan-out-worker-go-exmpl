package data_access

import "time"

type Post struct {
	ID        uint      `db:"id"`
	SenderID  uint      `db:"sender_id"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`
}

type SubscriberUser struct {
	SenderID     uint `db:"sender_id"`
	SubscriberID uint `db:"subscriber_id"`
}
