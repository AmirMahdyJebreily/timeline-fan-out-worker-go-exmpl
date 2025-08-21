package dataaccess

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type DataAccess struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *DataAccess {
	return &DataAccess{db: db}
}

func (da *DataAccess) InsertPost(ctx context.Context, post Post) (uint, time.Time, error) {
	tx, err := da.db.BeginTxx(ctx, nil)
	if err != nil {
		return 0, time.Time{}, fmt.Errorf(errBeginTx, err)
	}

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()

	insertQuery := tx.Rebind(`
        INSERT INTO posts (sender_id, content)
        VALUES (:sender_id, :content)
    `)

	res, err := tx.NamedExecContext(ctx, insertQuery, post)
	if err != nil {
		return 0, time.Time{}, fmt.Errorf(errInsertPosts, err)
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, time.Time{}, fmt.Errorf(errGetInsertedPostId, err)
	}
	insertedID := uint(lastID)

	var createdAt time.Time
	selectQuery := tx.Rebind(`SELECT created_at FROM posts WHERE id = ?`)
	if err := tx.GetContext(ctx, &createdAt, selectQuery, insertedID); err != nil {
		return 0, time.Time{}, fmt.Errorf("failed to fetch created_at for post %d: %w", insertedID, err)
	}

	if err := tx.Commit(); err != nil {
		return 0, time.Time{}, fmt.Errorf(errCommitTx, err)
	}
	committed = true

	return insertedID, createdAt, nil
}

func (da *DataAccess) GetSubscribers(ctx context.Context, userID uint) ([]uint, error) {
	query := `SELECT subscriber_id FROM subscriber_users WHERE sender_id = ?`
	tx, err := da.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf(errBeginTx, err)
	}

	query = tx.Rebind(query)

	var subscribers []uint
	err = tx.SelectContext(ctx, &subscribers, query, userID)
	if err != nil {
		RollbackOnError(tx)
		return nil, fmt.Errorf(errSelectSubscribers, err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf(errCommitTx, err)
	}

	return subscribers, nil
}

func (da *DataAccess) BulkGetPosts(ctx context.Context, postIDs []uint) ([]Post, error) {
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
		return nil, fmt.Errorf(errConstructPostsQuery, err)
	}

	tx, err := da.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf(errBeginTx, err)
	}

	var posts []Post

	err = tx.SelectContext(ctx, &posts, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return []Post{}, nil
		}
		defer RollbackOnError(tx)
		return nil, fmt.Errorf(errSelectPosts, err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf(errCommitTx, err)
	}

	return posts, nil
}
