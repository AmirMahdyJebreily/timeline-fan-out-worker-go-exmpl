package database

import (
	"context"
	"database/sql"

	"github.com/AmirMahdyJebreily/timeline-example/internal/database/entities"
)

type PostsRepository interface {
	InsertPost(ctx context.Context, posts []entities.Post) (resp sql.Result, err error)
	GetPostByIds(ctx context.Context, postIds []int) (posts []entities.Post, err error)
}
