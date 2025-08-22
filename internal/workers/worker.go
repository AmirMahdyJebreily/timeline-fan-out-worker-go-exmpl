package workers

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

type Job func() error

type WorkerPool struct {
	minWorkers   uint
	maxWorkers   uint
	ctx          context.Context
	cancel       context.CancelFunc
	jobs         chan Job
	errCh        chan error
	wg           sync.WaitGroup
	shutdownOnce sync.Once
}

var (
	instance WorkerPool
	once     sync.Once
)

func New(ctx context.Context, maximumWorkersCount, minimumWorkersCount, cap uint) (*WorkerPool, error) {
	if maximumWorkersCount < minimumWorkersCount {
		return nil, errors.New("maximum workers cannot be less than minimum workers")
	}
	once.Do(func() {
		cctx, cancel := context.WithCancel(ctx)
		instance = WorkerPool{
			ctx:        cctx,
			cancel:     cancel,
			minWorkers: minimumWorkersCount,
			maxWorkers: maximumWorkersCount,
			jobs:       make(chan Job, cap),
			errCh:      make(chan error, int(maximumWorkersCount)),
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
				err := job()
				if err != nil {
					select {
					case wp.errCh <- fmt.Errorf("worker %d: %w", id, err):
					default:
					}
				}
			}()
		}
	}
}

func (wp *WorkerPool) SubmitJob(j Job) error {
	select {
	case <-wp.ctx.Done():
		return wp.ctx.Err()
	case wp.jobs <- j:
		return nil
	}
}

func (wp *WorkerPool) Shutdown() {
	wp.shutdownOnce.Do(func() {
		wp.cancel()
		close(wp.jobs)
		wp.wg.Wait()
		close(wp.errCh)
	})
}
