package dataaccess

import (
	"context"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func setupTestDB(t *testing.T) *sqlx.DB {
	dsn := "admin:admin@tcp(localhost:3306)/timeline_db?parseTime=true&multiStatements=true"
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		t.Fatalf("failed to connect to test db: %v", err)
	}
	return db
}

func TestAddUserAndSubscriber(t *testing.T) {
	db := setupTestDB(t)
	da := New(db)
	ctx := context.Background()

	// Clean up tables before test
	_, _ = db.Exec("DELETE FROM subscriber_users")
	_, _ = db.Exec("DELETE FROM posts")
	_, _ = db.Exec("DELETE FROM users")

	// Add users
	res, err := db.Exec("INSERT INTO users (id) VALUES (?), (?)", 101, 202)
	if err != nil {
		t.Fatalf("failed to insert users: %v", err)
	}
	rows, _ := res.RowsAffected()
	if rows != 2 {
		t.Fatalf("expected 2 users inserted, got %d", rows)
	}

	// Add subscriber relationship: 202 subscribes to 101
	_, err = db.Exec("INSERT INTO subscriber_users (sender_id, subscriber_id) VALUES (?, ?)", 101, 202)
	if err != nil {
		t.Fatalf("failed to insert subscriber: %v", err)
	}

	// Test retrieval
	subs, err := da.BulkGetSubscribers(ctx, []uint{101})
	if err != nil {
		t.Fatalf("BulkGetSubscribers failed: %v", err)
	}
	if len(subs) != 1 {
		t.Fatalf("expected 1 subscriber, got %d", len(subs))
	}
	if subs[0].SenderID != 101 || subs[0].SubscriberID != 202 {
		t.Fatalf("unexpected subscriber data: %+v", subs[0])
	}
}

func TestBulkInsertPostsAndGetPosts(t *testing.T) {
	db := setupTestDB(t)
	da := New(db)
	ctx := context.Background()

	posts := []Post{}
	max := 1000
	for i := 0; i < max; i++ {
		posts = append(posts, Post{SenderID: 101, Content: fmt.Sprintf("Content No.%v", i)})
	}
	ids, err := da.BulkInsertPosts(ctx, posts)
	if err != nil {
		t.Fatalf("BulkInsertPosts failed: %v", err)
	}
	if len(ids) != max {
		t.Fatalf("expected %d ids, got %v", max, ids)
	}

	fetched, err := da.BulkGetPosts(ctx, ids)
	if err != nil {
		t.Fatalf("BulkGetPosts failed: %v", err)
	}
	if len(fetched) != max {
		t.Fatalf("expected %d posts, got %d", max, len(fetched))
	}
}

func TestBulkGetSubscribers_Empty(t *testing.T) {
	db := setupTestDB(t)
	da := New(db)
	ctx := context.Background()

	subs, err := da.BulkGetSubscribers(ctx, []uint{999999})
	if err != nil {
		t.Fatalf("BulkGetSubscribers failed: %v", err)
	}
	if len(subs) != 0 {
		t.Fatalf("expected 0 subscribers, got %d", len(subs))
	}
}
