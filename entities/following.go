package entities

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

const (
	HOME_TIMELINE_SIZE = 1000
)

func FollowUser(client *redis.Client, uid string, otherUid string) error {
	fkey1 := fmt.Sprintf("following:%s", uid)
	fkey2 := fmt.Sprintf("followers:%s", otherUid)

	if _, err := client.ZScore(fkey1, otherUid).Result(); err == nil {
		return errors.New(fmt.Sprintf("%s already followed %s", uid, otherUid))
	}

	now := time.Now().UnixNano() / int64(time.Millisecond)
	pipe := client.Pipeline()
	pipe.ZAdd(fkey1, redis.Z{
		Score: float64(now),
		Member: otherUid,
	})
	pipe.ZAdd(fkey2, redis.Z{
		Score: float64(now),
		Member: uid,
	})
	following := pipe.ZCard(fkey1)
	followers := pipe.ZCard(fkey2)
	statusAndSCores := pipe.ZRevRangeWithScores(fmt.Sprintf("profile:%s", otherUid),
		0, HOME_TIMELINE_SIZE - 1)
	_, err := pipe.Exec()
	fmt.Println(followers.Val(), err)

	pipe = client.Pipeline()

	pipe.HSet(fmt.Sprintf("users:%s", uid), "following", following.Val())
	pipe.HSet(fmt.Sprintf("users:%s", otherUid), "followers", followers.Val())
	if len(statusAndSCores.Val()) > 0 {
		pipe.ZAdd(fmt.Sprintf("home:%s", uid), statusAndSCores.Val()...)
	}

	pipe.ZRemRangeByRank(fmt.Sprintf("home:%s", uid), 0, -(HOME_TIMELINE_SIZE - 1))
	pipe.Exec()

	return nil
}
