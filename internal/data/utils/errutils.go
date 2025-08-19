package utils

import (
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"
)

const (
	ErrBeginTx         = "failed to begin transaction with db: %w"
	ErrBulkInsertPosts = "failed to bulk insert posts: %w"
	ErrLastInsertID    = "failed to get the last insert id: %w"
	ErrCommitTx        = "failed to commit transaction: %w"

	ErrConstructSubscribersQuery = "failed to construct query for getting subscribers: %w"
	ErrSelectSubscribers         = "failed to execute query for getting subscribers: %w"

	ErrConstructPostsQuery = "failed to construct query for getting posts: %w"
	ErrSelectPosts         = "failed to execute query for getting posts: %w"
)

func RollbackOnError(tx *sqlx.Tx) func() {
	return func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			log.Printf("rollback error: %v", err)
		}
	}
}
