package cache

import "fmt"

func timelineKey(userID uint) string {
	return fmt.Sprintf("timeline:%d", userID)
}
