package cache

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type TimelineCache interface {
	AddPostToTimeline(ctx context.Context, userID uint, postID uint, score float64) error
	GetTimelinePostIDs(ctx context.Context, userID uint, start, stop int64) ([]uint, error)
}

type timelineCache struct {
	client *redis.Client
}

func New(client *redis.Client) TimelineCache {
	return &timelineCache{client: client}
}

func (tc *timelineCache) AddPostToTimeline(ctx context.Context, userID uint, postID uint, score float64) error {
	key := timelineKey(userID)
	return tc.client.ZAdd(ctx, key, redis.Z{Score: score, Member: postID}).Err()
}

func (tc *timelineCache) GetTimelinePostIDs(ctx context.Context, userID uint, start, stop int64) ([]uint, error) {
	key := timelineKey(userID)
	ids, err := tc.client.ZRevRange(ctx, key, start, stop).Result()
	if err != nil {
		return nil, err
	}
	result := make([]uint, len(ids))
	for i, idStr := range ids {
		var id uint
		_, err := fmt.Sscan(idStr, &id)
		if err != nil {
			return nil, err
		}
		result[i] = id
	}
	return result, nil
}

func NewRedisClient(addr string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   0,
	})
}

func timelineKey(userID uint) string {
	return fmt.Sprintf("timeline:%d", userID)
}
