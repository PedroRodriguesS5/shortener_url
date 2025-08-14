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
	db, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = os.Getenv("REDIS_URL")
	}

	// Handle full Redis connection string (production)
	if strings.HasPrefix(redisAddr, "redis://") {
		opt, err := redis.ParseURL(redisAddr)
		if err != nil {
			log.Fatalf("error parsing Redis URL: %v", err)
		}
		Rdb = redis.NewClient(opt)
	} else {
		log.Printf("Connecting to Redis at: %s", redisAddr)
		Rdb = redis.NewClient(&redis.Options{
			Addr:     redisAddr,
			Username: os.Getenv("REDIS_USERNAME"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       db,
		})
	}

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

// GetAllURLs returns all shortened URLs with their original URLs and click counts
func GetAllURLs() (map[string]map[string]interface{}, error) {
	result := make(map[string]map[string]interface{})

	// Get all keys that are not click counters
	keys, err := Rdb.Keys(Ctx, "*").Result()
	if err != nil {
		return nil, err
	}

	for _, key := range keys {
		// Skip click counter keys
		if strings.HasPrefix(key, "clicks:") || strings.HasPrefix(key, "url_clicks:") {
			continue
		}

		// Get the original URL
		originalURL, err := Rdb.Get(Ctx, key).Result()
		if err != nil {
			continue // Skip if error getting URL
		}

		// Get click count (default to 0 if not found)
		clicks, err := Rdb.Get(Ctx, "clicks:"+key).Int64()
		if err != nil {
			clicks = 0
		}

		result[key] = map[string]interface{}{
			"original_url": originalURL,
			"clicks":       clicks,
		}
	}

	return result, nil
}
