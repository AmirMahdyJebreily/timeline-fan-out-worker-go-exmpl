package cache_test

import (
	"context"
	"testing"
	"time"

	"github.com/AmirMahdyJebreily/timeline-example/internal/cache"
	"github.com/redis/go-redis/v9"
)

func newTestCache(t *testing.T) *cache.TimelineCache {
	t.Helper()
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   9,
	})

	if err := client.FlushDB(context.Background()).Err(); err != nil {
		t.Fatalf("failed to flush redis: %v", err)
	}

	return cache.New(client)
}

func TestAddPostToTimeline(t *testing.T) {
	ctx := context.Background()
	tc := newTestCache(t)

	now := float64(time.Now().UnixMilli())
	if err := tc.AddPostToTimeline(ctx, 1, 42, now); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_ = tc.AddPostToTimeline(ctx, 1, 100, now+1000)
	_ = tc.AddPostToTimeline(ctx, 1, 200, now+500)

	ids, err := tc.GetTimelinePostIDs(ctx, 1, 0, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []uint{100, 200, 42}
	if len(ids) != len(expected) {
		t.Fatalf("expected %d ids, got %d", len(expected), len(ids))
	}
	for i := range expected {
		if ids[i] != expected[i] {
			t.Fatalf("expected %v at %d, got %v", expected[i], i, ids[i])
		}
	}
}

func TestGetTimelinePostIDs(t *testing.T) {
	ctx := context.Background()
	tc := newTestCache(t)

	ids, err := tc.GetTimelinePostIDs(ctx, 2, 0, 10)
	if err != nil {
		t.Fatalf("unexpected error on empty timeline: %v", err)
	}
	if len(ids) != 0 {
		t.Fatalf("expected empty slice, got %v", ids)
	}

	now := float64(time.Now().UnixMilli())
	_ = tc.AddPostToTimeline(ctx, 2, 1, now)
	_ = tc.AddPostToTimeline(ctx, 2, 2, now+10)
	_ = tc.AddPostToTimeline(ctx, 2, 3, now+20)

	ids, err = tc.GetTimelinePostIDs(ctx, 2, 0, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []uint{3, 2}
	if len(ids) != len(expected) {
		t.Fatalf("expected %d ids, got %d", len(expected), len(ids))
	}
	for i := range expected {
		if ids[i] != expected[i] {
			t.Fatalf("expected %v at %d, got %v", expected[i], i, ids[i])
		}
	}
}
