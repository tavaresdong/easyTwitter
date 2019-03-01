package entities

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"strconv"
	"strings"
	"time"
)

func Double (i int) int {
	return i * 2
}

func DeleteUser(client *redis.Client, login string) error {
	login = strings.ToLower(login)
	id, err := client.HGet("users:", login).Result()
	if err == redis.Nil {
		// If the login name already exist, return false
		return errors.New("The login key does not exist")
	} else if err != nil {
		return err
	}

	lock := AcquireLockWithTimeout(client, "user:" + login, 10, 10)
	if lock == "" {
		fmt.Println("User metadata is modified by another client, please wait")
		return errors.New("Failed to acquire lock")
	}
	defer ReleaseLock(client, "user:" + login, lock)

	ndelUser, e1 := client.HDel("users:", login).Result()
	ndelInfo, e2 := client.Del("users:" + id).Result()

	if e1 != nil {
		return e1
	}

	if e2 != nil {
		return e2
	}

	if ndelUser != 1 || ndelInfo != 1 {
		return errors.New("Number of elements deleted is not 1")
	}

	return nil
}

func FetchUser(client *redis.Client, login string) map[string]interface{} {
	login = strings.ToLower(login)
	id, err := client.HGet("users:", login).Result()
	if err == redis.Nil {
		fmt.Printf("Key %s does not exist!", login)
		return nil
	} else if err != nil {
		panic(err)
	}

	results, err := client.HMGet("users:" + id,
		"login", "id", "name", "followers", "following", "posts", "signup").Result()
	if err != nil {
		panic(err)
	}

	m := map[string]interface{} {
		"login": results[0],
		"id": results[1],
		"name": results[2],
		"followers": results[3],
		"following": results[4],
		"posts": results[5],
		"signup": results[6],
	}

	return m
}

func CreateUser(client *redis.Client, login string, name string) string {
	login = strings.ToLower(login)
	lock := AcquireLockWithTimeout(client, "user:" + login, 10, 10)
	if lock == "" {
		return ""
	}
	defer ReleaseLock(client, "user:" + login, lock)

	hget := client.HGet("users:", login)
	if hget.Err() == nil || hget.Err() != redis.Nil {
		// If the login name already exist, return empty
		return ""
	}

	id, err := client.Incr("user:id:").Result()
	if err != nil {
		panic(err)
	}

	pipe := client.Pipeline()
	pipe.HSet("users:", login, id)

	m := map[string]interface{} {
		"login":login,
		"id":id,
		"name":name,
		"followers":0,
		"following":0,
		"posts":0,
		"signup":time.Now().UnixNano() / int64(time.Millisecond),
	}

	uid := strconv.FormatInt(id, 10)
	result, err := pipe.HMSet("users:" + uid, m).Result()
	pipe.Exec()

	if err != nil {
		panic(err)
	}

	fmt.Println(result)
	return strconv.FormatInt(id, 10)
}
