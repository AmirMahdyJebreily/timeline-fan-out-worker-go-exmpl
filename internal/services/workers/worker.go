package workers

import (
	"context"
	"time"
	"log"

	"github.com/AmirMahdyJebreily/timeline-example/internal/cache"
	dataaccess "github.com/AmirMahdyJebreily/timeline-example/internal/data"
)

type fanOutWorker struct {
	db    dataaccess.DataAccesser
	cache cache.TimelineCache
}

type FanOutWorker interface {
	AddWorker(id int, ctx context.Context, jobs <-chan WorkerJob, results chan<- dataaccess.Post)
}

func New(db dataaccess.DataAccesser, cache cache.TimelineCache) *fanOutWorker {
	return &fanOutWorker{db: db, cache: cache}
}

func (supply *fanOutWorker) AddWorker(id int, ctx context.Context, jobs <-chan WorkerJob, results chan<- dataaccess.Post) {
	for job := range jobs {
			// Use Redis pipeline for batch efficiency (if you want to extend to batch jobs)
			// For now, single ZADD per job (per subscriber)
			err := supply.cache.AddPostToTimeline(ctx, job.SubscriberID, job.PostID, float64(time.Now().UnixNano()))
			if err != nil {
				log.Printf("worker %d: failed to add post %d to subscriber %d timeline: %v", id, job.PostID, job.SubscriberID, err)
				continue
			}
			// Optionally, fetch the post and send to results channel for further processing/ack
			if results != nil {
				posts, err := supply.db.BulkGetPosts(ctx, []uint{job.PostID})
				if err == nil && len(posts) > 0 {
					results <- posts[0]
				}
			}
		}
}
