package dataaccess

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/AmirMahdyJebreily/timeline-example/internal/data/utils"
	"github.com/jmoiron/sqlx"
)

type dataAccess struct {
	db *sqlx.DB
}

type DataAccesser interface {
	BulkInsertPosts(ctx context.Context, posts []Post) ([]uint, error)
	BulkGetSubscribers(ctx context.Context, userIDs []uint) ([]SubscriberUser, error)
	BulkGetPosts(ctx context.Context, postIDs []uint) ([]Post, error)
}

func New(db *sqlx.DB) DataAccesser {
	return &dataAccess{db: db}
}

func (da *dataAccess) BulkInsertPosts(ctx context.Context, posts []Post) ([]uint, error) {
	if len(posts) == 0 {
		return nil, nil
	}

	tx, err := da.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf(utils.ErrBeginTx, err)
	}

	valueStrings := make([]string, 0, len(posts))
	valueArgs := make([]interface{}, 0, len(posts)*2)
	for _, post := range posts {
		valueStrings = append(valueStrings, "(?, ?)")
		valueArgs = append(valueArgs, post.SenderID, post.Content)
	}
	query := fmt.Sprintf("INSERT INTO posts (sender_id, content) VALUES %s", strings.Join(valueStrings, ","))

	res, err := tx.ExecContext(ctx, query, valueArgs...)
	if err != nil {
		defer utils.RollbackOnError(tx)
		return nil, fmt.Errorf(utils.ErrBulkInsertPosts, err)
	}

	firstID, err := res.LastInsertId()
	if err != nil {
		defer utils.RollbackOnError(tx)
		return nil, fmt.Errorf(utils.ErrLastInsertID, err)
	}

	ids := make([]uint, len(posts))
	for i := range posts {
		ids[i] = uint(firstID) + uint(i)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf(utils.ErrCommitTx, err)
	}
	return ids, nil
}

func (da *dataAccess) BulkGetSubscribers(ctx context.Context, userIDs []uint) ([]SubscriberUser, error) {
	if len(userIDs) == 0 {
		return nil, nil
	}

	query := `SELECT sender_id, subscriber_id FROM subscriber_users WHERE sender_id IN (?)`

	query, args, err := sqlx.In(query, userIDs)
	if err != nil {
		return nil, fmt.Errorf(utils.ErrConstructSubscribersQuery, err)
	}

	tx, err := da.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf(utils.ErrBeginTx, err)
	}

	query = tx.Rebind(query)

	var subscribers []SubscriberUser
	err = tx.SelectContext(ctx, &subscribers, query, args...)
	if err != nil {
		defer utils.RollbackOnError(tx)
		return nil, fmt.Errorf(utils.ErrSelectSubscribers, err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf(utils.ErrCommitTx, err)
	}

	return subscribers, nil
}

func (da *dataAccess) BulkGetPosts(ctx context.Context, postIDs []uint) ([]Post, error) {
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
		return nil, fmt.Errorf(utils.ErrConstructPostsQuery, err)
	}

	tx, err := da.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf(utils.ErrBeginTx, err)
	}

	query = tx.Rebind(query)

	var posts []Post

	err = tx.SelectContext(ctx, &posts, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return []Post{}, nil
		}
		defer utils.RollbackOnError(tx)
		return nil, fmt.Errorf(utils.ErrSelectPosts, err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf(utils.ErrCommitTx, err)
	}

	return posts, nil
}
