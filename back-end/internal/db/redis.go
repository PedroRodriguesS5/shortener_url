package db

import (
	"context"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()
var Rdb *redis.Client

func InitRedis() {
	db, _ := strconv.Atoi(os.Getenv("DB"))

	// Check for REDIS_ADDR first (used in Render and docker-compose)
	// Fall back to REDIS_URL if REDIS_ADDR is not set
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = os.Getenv("REDIS_URL")
	}

	// Handle redis:// URL format from Render
	// Extract only the host:port part
	if strings.HasPrefix(redisAddr, "redis://") {
		// Remove the redis:// prefix
		redisAddr = strings.TrimPrefix(redisAddr, "redis://")
		// Remove any path or query parameters (everything after a slash)
		if slashIndex := strings.Index(redisAddr, "/"); slashIndex != -1 {
			redisAddr = redisAddr[:slashIndex]
		}
		// Remove any auth info (username:password@)
		if atIndex := strings.Index(redisAddr, "@"); atIndex != -1 {
			redisAddr = redisAddr[atIndex+1:]
		}
	}

	log.Printf("Connecting to Redis at: %s", redisAddr)

	Rdb = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
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

func IncrementURLClick(fullURL string) error {
	return Rdb.Incr(Ctx, "url_clicks:"+fullURL).Err()
}

func GetClicks(code string) (int64, error) {
	return Rdb.Get(Ctx, "clicks:"+code).Int64()
}

func GetURLClicks(fullURL string) (int64, error) {
	return Rdb.Get(Ctx, "url_clicks:"+fullURL).Int64()
}
