package entities

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/satori/go.uuid"
	"strings"
	"time"
)

func ConnectToRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "yucdongredis.redis.cache.windows.net:6379",
		Password: "Sang8sMBUmDsvn8dIWVvd1n7DUvFE5DzPAa8E7tN0nc=", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

	return client
}

func ClearLock(client *redis.Client, lockname string) bool {
	lockname = "lock:" + lockname
	val, _ := client.Del(lockname).Result()
	return val == 1
}

func AcquireLockWithTimeout(client *redis.Client, lockname string, acquireTimeout int64, lockTimeout int64) string {
	identifier, _ := uuid.NewV4()

	lockname = "lock:" + lockname

	ending := time.Now().Add(time.Second * time.Duration(acquireTimeout))
	for time.Now().Before(ending) {
		succ, _ := client.SetNX(lockname, identifier.String(), 0).Result()
		if succ {
			client.Expire(lockname, time.Duration(lockTimeout) * time.Second)
			return identifier.String()
		} else {
			ttl, _ := client.TTL(lockname).Result()
			if ttl < 0 {
				client.Expire(lockname, time.Duration(lockTimeout) * time.Second)
			}
		}

		time.Sleep(time.Millisecond * 10)
	}

	return ""
}

func AcquireLock(client *redis.Client, lockname string, lockTimeout int64) string {
	identifier, _ := uuid.NewV4()

	lockname = "lock:" + lockname

	succ, _ := client.SetNX(lockname, identifier.String(), 0).Result()
	if succ {
		client.Expire(lockname, time.Duration(lockTimeout) * time.Second)
		return identifier.String()
	} else {
		ttl, _ := client.TTL(lockname).Result()
		if ttl < 0 {
			client.Expire(lockname, time.Duration(lockTimeout) * time.Second)
		}
	}

	return ""
}

func ReleaseLock(client *redis.Client, lockname string, identifier string) error {
	var release func(string) error
	lockname = "lock:" + lockname

	// Transactionally increments key using GET and SET commands.
	release = func(key string) error {
		err := client.Watch(func(tx *redis.Tx) error {
			id, err := tx.Get(key).Result()
			if err != nil || strings.Compare(id, identifier) != 0 {
				return errors.New("Identifier mismatch")
				tx.Unwatch(key)
			}

			_, e := tx.Pipelined(func(pipe redis.Pipeliner) error {
				pipe.Del(lockname)
				return nil
			})

			return e
		}, lockname)
		if err == redis.TxFailedErr {
			return release(key)
		}
		return err
	}

	return release(lockname)
}