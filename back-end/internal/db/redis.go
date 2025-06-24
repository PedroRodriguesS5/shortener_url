package db

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()
var Rdb *redis.Client

func InitRedis() {
	db, _ := strconv.Atoi(os.Getenv("DB"))
	Rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
	})

	_, err := Rdb.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("error connecting to Redis: %v", err)
	}
}

func SaveURL(code, url string, expired int) error {
	return Rdb.Set(Ctx, code, url, time.Duration(expired)*time.Second).Err()
}

func GetURL(code string) (string, error) {
	return Rdb.Get(Ctx, code).Result()
}

func IncrementClick(code string) error {
	return Rdb.Incr(Ctx, "clicks:"+code).Err()
}

func GetClicks(code string) (int64, error) {
	return Rdb.Get(Ctx, "clicks:"+code).Int64()
}
