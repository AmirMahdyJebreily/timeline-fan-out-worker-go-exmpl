package timeline

import (
	"context"
	"fmt"
	"log"
	"sync"
)

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
