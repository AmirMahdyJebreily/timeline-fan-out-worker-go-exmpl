package data_access

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type dataAccess struct {
	db *sqlx.DB
}

type DataAccesser interface {
	BulkInsertPosts(posts []Post) error
	BulkGetSubscribers(userIDs []uint) ([]SubscriberUser, error)
	BulkGetPosts(postIDs []uint) ([]Post, error)
}

func New(db *sqlx.DB) DataAccesser {
	return &dataAccess{db: db}
}

func (da *dataAccess) BulkInsertPosts(posts []Post) error {
	if len(posts) == 0 {
		return nil
	}

	query := `INSERT INTO posts (sender_id, content) VALUES (:sender_id, :content)`

	_, err := da.db.NamedExec(query, posts)
	if err != nil {
		return fmt.Errorf("failed to bulk insert posts: %w", err)
	}

	return nil
}

func (da *dataAccess) BulkGetSubscribers(userIDs []uint) ([]SubscriberUser, error) {
	if len(userIDs) == 0 {
		return nil, nil
	}

	query := `SELECT sender_id, subscriber_id FROM subscriber_users WHERE sender_id IN (?)`

	query, args, err := sqlx.In(query, userIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to construct query for getting subscribers: %w", err)
	}

	query = da.db.Rebind(query)

	var subscribers []SubscriberUser
	err = da.db.Select(&subscribers, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query for getting subscribers: %w", err)
	}

	return subscribers, nil
}

func (da *dataAccess) BulkGetPosts(postIDs []uint) ([]Post, error) {
	if len(postIDs) == 0 {
		return nil, nil
	}

	query := `SELECT id, sender_id, content, created_at FROM posts WHERE id IN (?) ORDER BY FIND_IN_SET(id, ?)`

	idStrings := make([]string, len(postIDs))
	for i, id := range postIDs {
		idStrings[i] = fmt.Sprint(id)
	}
	orderedIDList := strings.Join(idStrings, ",")

	query, args, err := sqlx.In(query, postIDs, orderedIDList)
	if err != nil {
		return nil, fmt.Errorf("failed to construct query for getting posts: %w", err)
	}

	query = da.db.Rebind(query)

	var posts []Post

	err = da.db.Select(&posts, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return []Post{}, nil
		}
		return nil, fmt.Errorf("failed to execute query for getting posts: %w", err)
	}

	return posts, nil
}
