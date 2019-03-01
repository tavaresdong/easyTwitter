package entities

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

func FetchStatus(client *redis.Client, statusid int64) (map[string]interface{}, error) {
	result, err := client.HMGet(fmt.Sprintf("status:%d", statusid),
		"message", "posted", "id", "uid", "login").Result()

	if err == redis.Nil {
		return nil, errors.New(fmt.Sprintf("No such posts exists: %d", statusid))
	} else if err != nil {
		return nil, err
	} else {
		data := map[string]interface{} {
			"message": result[0],
			"posted": result[1],
			"id": result[2],
			"uid": result[3],
			"login": result[4],
		}

		return data, nil
	}
}


func CreateStatus(client *redis.Client, uid string, message string) (int64, error) {
	login, err := client.HGet("users:" + uid, "login").Result()
	if err == redis.Nil {
		return -1, errors.New("No such user: " + uid)
	}

	id, err := client.Incr("status:id:").Result()
	if err != nil {
		return -1, err
	}

	pipe := client.Pipeline()
	data := map[string]interface{} {
		"message": message,
		"posted": time.Now().UnixNano() / int64(time.Millisecond),
		"id": id,
		"uid": uid,
		"login": login,
	}
	pipe.HMSet(fmt.Sprintf("status:%d", id), data)
	pipe.HIncrBy(fmt.Sprintf("user:%s", uid), "posts", 1)
	pipe.Exec()

	return id, nil
}