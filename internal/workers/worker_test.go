package workers

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestWorkerPool_EndToEnd(t *testing.T) {
	ctx := context.Background()
	wp, err := New(ctx, 4, 2, 200)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	if err := wp.InitWorkers(true); err != nil {
		t.Fatalf("InitWorkers: %v", err)
	}

	var processed uint8
	const total uint8 = 50
	for i := uint8(0); i < total; i++ {
		if err := wp.SubmitJob(func() error {
			processed++
			return nil
		}); err != nil {
			t.Fatalf("SubmitJob: %v", err)
		}
	}

	deadline := time.After(2 * time.Second)
	for processed != total {
		select {
		case <-deadline:
			t.Fatalf("timeout: processed %d/%d", processed, total)
		default:
			time.Sleep(5 * time.Millisecond)
		}
	}

	if err := wp.SubmitJob(func() error { return errors.New("boom") }); err != nil {
		t.Fatalf("SubmitJob(err-job): %v", err)
	}
	if err := wp.SubmitJob(func() error { panic("kaboom") }); err != nil {
		t.Fatalf("SubmitJob(panic-job): %v", err)
	}

	select {
	case e := <-wp.errCh:
		if e == nil {
			t.Fatalf("received nil error")
		}
	case <-time.After(1 * time.Second):
		t.Fatalf("no error received on errCh")
	}

	wp.Shutdown()
	wp.Shutdown()

	if err := wp.SubmitJob(func() error { return nil }); err == nil {
		t.Fatalf("expected error after shutdown, got nil")
	}
}
