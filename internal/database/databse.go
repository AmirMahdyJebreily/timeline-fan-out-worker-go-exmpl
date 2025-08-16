package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/AmirMahdyJebreily/timeline-example/internal/database/entities"
)

type mysqlPostsRepository struct {
	db *sql.DB
}

func New(db *sql.DB) (resp PostsRepository) {
	return &mysqlPostsRepository{
		db,
	}
}

func (repo *mysqlPostsRepository) InsertPost(ctx context.Context, posts []entities.Post) (resp sql.Result, err error) {
	jsonData, err := json.Marshal(posts)
	if err != nil {
		err = fmt.Errorf("Parsing posts: %v", err)
		log.Println(err)
		return
	}

	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		err = fmt.Errorf("Bgining transaction: %v", err)
		log.Println(err)
		return
	}

	defer func() {
		if err != nil {
			rbErr := tx.Rollback()
			err = fmt.Errorf("Executing transaction: %v\n Rollback: %v", err, rbErr)
			log.Println(err)
			return
		}
	}()

	resp, err = tx.ExecContext(ctx, "CALL sp_insert_posts_bulk(?)", string(jsonData))

	if err != nil {
		err = fmt.Errorf("Parsing posts: %v", err)
		log.Println(err)
		return
	}

	if err = tx.Commit(); err != nil {
		err = fmt.Errorf("Committing transaction: %v", err)
		log.Println(err)
		return
	}

	return
}

func (repo *mysqlPostsRepository) GetPostByIds(ctx context.Context, postIds []int) (posts []entities.Post, err error) {
	return
}
