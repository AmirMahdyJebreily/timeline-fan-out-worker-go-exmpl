package timeline

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/AmirMahdyJebreily/timeline-example/internal/cache"
	dataaccess "github.com/AmirMahdyJebreily/timeline-example/internal/data"
	"github.com/AmirMahdyJebreily/timeline-example/internal/workers"
)

type TimelineService struct {
	db      *dataaccess.DataAccess
	cache   *cache.TimelineCache
	workers *workers.WorkerPool
}

var (
	instance TimelineService
	once     sync.Once
)

func New(db *dataaccess.DataAccess, cache *cache.TimelineCache, workers *workers.WorkerPool) *TimelineService {
	once.Do(func() {
		instance = TimelineService{
			db:      db,
			cache:   cache,
			workers: workers,
		}
	})
	return &instance
}

func (tl *TimelineService) NewPost(ctx context.Context, post dataaccess.Post) (uint, error) {
	id, at, err := tl.db.InsertPost(ctx, post)
	if err != nil {
		return 0, err
	}

	select {
	case <-ctx.Done():
		return id, ctx.Err()
	default:
	}

	go tl.PostToSubs(context.Background(), post.SenderID, id, at)

	return id, nil
}

func (tl *TimelineService) PostToSubs(ctx context.Context, userId, postId uint, at time.Time) error {
	subscribers, err := tl.db.GetSubscribers(ctx, userId)
	if err != nil {
		return fmt.Errorf("fail to get subscribers from sender id %w", err)
	}
	if len(subscribers) == 0 {
		return nil
	}

	score := float64(at.UnixMicro())

sendLoop:
	for _, sub := range subscribers {
		select {
		case <-ctx.Done():
			break sendLoop
		default:
			tl.workers.SubmitJob(tl.fanout(ctx, sub, postId, score))
		}
	}

	return nil
}

func (tl *TimelineService) fanout(ctx context.Context, userID, postID uint, score float64) workers.Job {
	return func() error {
		if err := tl.cache.AddPostToTimeline(ctx, userID, postID, score); err != nil {
			return err
		}
		return nil
	}
}
