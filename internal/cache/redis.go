package cache

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type TimelineCache struct {
	client *redis.Client
}

func New(client *redis.Client) *TimelineCache {
	return &TimelineCache{client: client}
}

func (tc *TimelineCache) AddPostToTimeline(ctx context.Context, userID uint, postID uint, score float64) error {
	key := timelineKey(userID)
	err := tc.client.ZAdd(ctx, key, redis.Z{Score: score, Member: postID}).Err()
	if err != nil {
		err = fmt.Errorf("fail to add to ZSET %w", err)
	}

	return err
}

func (tc *TimelineCache) GetTimelinePostIDs(ctx context.Context, userID uint, start, stop int64) ([]uint, error) {
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
