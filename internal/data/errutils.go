package dataaccess

import (
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"
)

const (
	errBeginTx      = "failed to begin transaction with db: %w"
	errInsertPosts  = "failed to insert posts: %w"
	errLastInsertID = "failed to get the last insertion id: %w"
	errCommitTx     = "failed to commit transaction: %w"

	errConstructPostsQuery = "failed to construct query for getting posts: %w"
	errSelectPosts         = "failed to execute query for getting posts: %w"
	errGetInsertedPostId   = "failed to get inserted post id: %w"
	errSelectSubscribers   = "failed to execute query for getting subscribers: %w"
	errGetSubscribersRows  = "failed to get subscribers query result in rows: %w"
	errGetSubscribers      = "failed to get subscribers query result: %w"
)

func RollbackOnError(tx *sqlx.Tx) func() {
	return func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			log.Printf("rollback error: %v", err)
		}
	}
}
