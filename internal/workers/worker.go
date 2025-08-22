package workers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
)

type Job func() (bool, error)

type WorkerPool struct {
	minWorkers uint
	maxWorkers uint
	ctx        context.Context
	cancel     context.CancelFunc
	jobs       chan Job
	errCh      chan error
	wg         sync.WaitGroup
}

var (
	instance WorkerPool
	once     sync.Once
)

func New(ctx context.Context, maximumWorkersCount, minimumWorkersCount, cap uint) (*WorkerPool, error) {
	if maximumWorkersCount < minimumWorkersCount {
		return nil, errors.New("maximum workers is more than minimum workers")
	}
	once.Do(func() {
		cctx, cancel := context.WithCancel(ctx)
		instance = WorkerPool{
			ctx:        cctx,
			cancel:     cancel,
			minWorkers: minimumWorkersCount,
			maxWorkers: maximumWorkersCount,
			jobs:       make(chan Job, cap),
			errCh:      make(chan error, 1),
		}
	})
	return &instance, nil
}

func (wp *WorkerPool) InitWorkers(max bool) error {
	var count uint = wp.maxWorkers
	if !max {
		count = wp.minWorkers
	}
	for i := uint(0); i < count; i++ {
		wp.wg.Add(1)
		go func(id uint) {
			defer wp.wg.Done()
			wp.worker(id)
		}(i)
	}
	return nil
}

func (wp *WorkerPool) worker(id uint) {
	for {
		select {
		case <-wp.ctx.Done():
			return
		case job, ok := <-wp.jobs:
			if !ok {
				return
			}
			func() {
				defer func() {
					if r := recover(); r != nil {
						select {
						case wp.errCh <- fmt.Errorf("worker %d panic: %v", id, r):
						default:
						}
					}
				}()
				okJob, err := job()
				if err != nil {
					select {
					case wp.errCh <- fmt.Errorf("worker %d: %w", id, err):
					default:
					}
				}
				if !okJob {
					log.Printf("worker %d: job returned false", id)
				}
			}()
		}
	}
}

func (wp *WorkerPool) Submit(j Job) error {
	select {
	case <-wp.ctx.Done():
		return wp.ctx.Err()
	case wp.jobs <- j:
		return nil
	}
}

func (wp *WorkerPool) Shutdown() {
	wp.cancel()
	close(wp.jobs)
	wp.wg.Wait()
	close(wp.errCh)
}
