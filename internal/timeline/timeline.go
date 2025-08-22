package timeline

import (
	"context"
	"fmt"
	"log"
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

const defaultMaxConcurrency = 50

func New(db *dataaccess.DataAccess, cache *cache.TimelineCache) *TimelineService {
	once.Do(func() {
		instance = TimelineService{
			db:    db,
			cache: cache,
		}
	})
	return &instance
}

func (tl *TimelineService) NewPost(ctx context.Context, post dataaccess.Post) (uint, time.Time, error) {
	// errSubs := tl.PostToSubs(ctx, post.SenderID, id)
	// if errSubs != nil {
	// 	err = fmt.Errorf("%w", err, errSubs)
	// }
	return tl.db.InsertPost(ctx, post)
}

func (tl *TimelineService) PostToSubs(ctx context.Context, userId, postId uint, at time.Time) error {
	subscribers, err := tl.db.GetSubscribers(ctx, userId)
	if err != nil {
		return fmt.Errorf("fail to get subscribers from sender id %w", err)
	}
	if len(subscribers) == 0 {
		return nil
	}

	jobs := make(chan uint)
	var wg sync.WaitGroup

	errCh := make(chan error, 1)

	score := float64(at.UnixMicro())

	for w := 0; w < defaultMaxConcurrency; w++ {
		wg.Add(1)
		go tl.fanoutWorkerOnce(ctx, w, jobs, &wg, postId, score, errCh)
	}

sendLoop:
	for _, sub := range subscribers {
		select {
		case <-ctx.Done():
			break sendLoop
		case jobs <- sub:
		}
	}

	close(jobs)
	wg.Wait()

	select {
	case e := <-errCh:
		return e
	default:
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return nil
	}

	// for i, subscriber := range subscribers {

	// 	go func() {
	// 		err := tl.cache.AddPostToTimeline(ctx, subscriber, postId, float64(at.UnixMicro()))
	// 		if err != nil {
	// 			err = fmt.Errorf("fan-out goroutine (no.%d) failed to add post into timeline %w", i, err)
	// 			log.Println(err)
	// 		}
	// 	}()
	// }
	//return nil
}

func (tl *TimelineService) fanoutWorkerOnce(ctx context.Context, workerID int, jobs <-chan uint, wg *sync.WaitGroup, postID uint, score float64, errCh chan<- error) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case sub, ok := <-jobs:
			if !ok {
				return
			}
			if err := tl.cache.AddPostToTimeline(ctx, sub, postID, score); err != nil {
				log.Printf("fan-out worker %d: subscriber=%d add failed: %v", workerID, sub, err)
				select {
				case errCh <- fmt.Errorf("worker %d failed in add post timeline to subscriber=%d: %w", workerID, sub, err):
				default:
				}
				return
			}
		}
	}
}
