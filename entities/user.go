package entities

import (
	"fmt"
	"github.com/go-redis/redis"
	"strconv"
	"strings"
	"time"
)

func Double (i int) int {
	return i * 2
}

func CreateUser(client *redis.Client, login string, name string) string {
	login = strings.ToLower(login)
	lock := AcquireLockWithTimeout(client, "user:" + login, 10, 10)
	if lock == "" {
		return ""
	}
	defer ReleaseLock(client, "user:" + login, lock)

	hget := client.HGet("users:", login)
	if hget.Err() == redis.Nil {
		return ""
	}

	id, _ := client.Incr("user:id:").Result()
	pipe := client.Pipeline()
	pipe.HSet("users:", login, id)

	m := map[string]interface{} {
		"login":login,
		"id":id,
		"name":name,
		"followers":0,
		"following":0,
		"posts":0,
		"signup":time.Now(),
	}
	result := pipe.HMSet("users:" + strconv.FormatInt(id, 10), m)
	pipe.Exec()

	fmt.Println(result.Val())
	return strconv.FormatInt(id, 10)
}
