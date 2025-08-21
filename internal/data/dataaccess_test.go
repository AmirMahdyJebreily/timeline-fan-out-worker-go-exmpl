package dataaccess

import (
	"context"
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

	_, _ = db.Exec("DELETE FROM subscriber_users")
	_, _ = db.Exec("DELETE FROM posts")
	_, _ = db.Exec("DELETE FROM users")

	const (
		TEST_SENDER_ID    = 101
		TEST_SUB_ID       = 202
		USERS_COUNT       = 2
		SUBSCRIBERS_COUNT = 1
	)

	res, err := db.Exec("INSERT INTO users (id) VALUES (?), (?)", TEST_SENDER_ID, TEST_SUB_ID)
	if err != nil {
		t.Fatalf("failed to insert users: %v", err)
	}

	rows, _ := res.RowsAffected()
	if rows != USERS_COUNT {
		t.Fatalf("expected 2 users inserted, got %d", rows)
	}

	_, err = db.Exec("INSERT INTO subscriber_users (sender_id, subscriber_id) VALUES (?, ?)", TEST_SENDER_ID, TEST_SUB_ID)
	if err != nil {
		t.Fatalf("failed to insert subscriber: %v", err)
	}

	subs, err := da.GetSubscribers(ctx, TEST_SENDER_ID)
	if err != nil {
		t.Fatalf("BulkGetSubscribers failed: %v", err)
	}
	if len(subs) != SUBSCRIBERS_COUNT {
		t.Fatalf("expected 1 subscriber, got %d", len(subs))
	}
	if subs[0] != TEST_SUB_ID {
		t.Fatalf("unexpected subscriber data: %+v", subs[0])
	}
}

func TestBulkInsertPostsAndGetPosts(t *testing.T) {
	db := setupTestDB(t)
	da := New(db)
	ctx := context.Background()
	const (
		CONTENT = "Content"
		COUNT   = 1
	)

	post := Post{SenderID: 101, Content: CONTENT}

	id, _, err := da.InsertPost(ctx, post)
	if err != nil {
		t.Fatalf("BulkInsertPosts failed: %v", err)
	}

	fetched, err := da.BulkGetPosts(ctx, []uint{id})
	if err != nil {
		t.Fatalf("BulkGetPosts failed: %v", err)
	}
	if len(fetched) > COUNT {
		t.Fatalf("expected %d in content, got %d", COUNT, len(fetched))
	}
	if fetched[0].Content != CONTENT {
		t.Fatalf("expected %v in content, got %v", CONTENT, fetched[0].Content)
	}
}

func TestBulkGetSubscribers_Empty(t *testing.T) {
	db := setupTestDB(t)
	da := New(db)
	ctx := context.Background()

	subs, err := da.GetSubscribers(ctx, 999999)
	if err != nil {
		t.Fatalf("BulkGetSubscribers failed: %v", err)
	}
	if len(subs) != 0 {
		t.Fatalf("expected 0 subscribers, got %d", len(subs))
	}
}
