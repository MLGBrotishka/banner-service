package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"my_app/internal/models"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var rdb *redis.Client
var ttl time.Duration

func InitCache() {
	host := os.Getenv("CACHE_HOST")
	port := os.Getenv("CACHE_PORT")
	password := os.Getenv("CACHE_PASSWORD")
	var err error
	ttl, err = time.ParseDuration(os.Getenv("CACHE_TTL"))
	if err != nil {
		log.Fatal(err)
	}
	rdb = redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       0,
	})

	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatal(err)
	}
}

func CloseCache() {
	err := rdb.Close()
	if err != nil {
		log.Fatalf("Error closing Redis connection: %v", err)
	}
}

func GetBannerFromCache(featureId *int, tagId *int) (*models.BannerExpanded, error) {
	var banner models.BannerExpanded
	cacheKey := fmt.Sprintf("banner:%d:%d", *featureId, *tagId)
	result, err := rdb.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("no banner found")
	}
	err = json.Unmarshal([]byte(result), &banner)
	if err != nil {
		return nil, err
	}
	log.Println("Loaded from cache")
	return &banner, nil
}

func SaveBannerToCacheAsync(featureId *int, tagId *int, banner *models.BannerExpanded) {
	go func(featureId int, tagId int, banner models.BannerExpanded) {
		err := SaveBannerToCache(&featureId, &tagId, &banner)
		if err != nil {
			log.Printf("Failed to save banner to cache: %v", err)
		}
	}(*featureId, *tagId, *banner)
}

func SaveBannerToCache(featureId *int, tagId *int, banner *models.BannerExpanded) error {
	cacheKey := fmt.Sprintf("banner:%d:%d", *featureId, *tagId)
	bannerJson, err := json.Marshal(banner)
	if err != nil {
		return err
	}
	err = rdb.Set(ctx, cacheKey, bannerJson, ttl).Err()
	if err != nil {
		return err
	}
	log.Println("Saved to cache")
	return err
}
