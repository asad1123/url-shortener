package cache

import (
	"context"

	"github.com/go-redis/redis/v8"
)

func connectToCache() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "cache:6379",
		Password: "", // no pw set for now
		DB:       0,  // use default db
	})

	return client
}

func SaveShortenedUrlToCache(ctx context.Context, redirectUrl string, shortenedId string) error {
	client := connectToCache()

	_, err := client.Set(ctx, shortenedId, redirectUrl, 0).Result()

	return err
}

func GetRedirectUrlFromCache(ctx context.Context, shortenedId string) (string, error) {
	client := connectToCache()

	redirectUrl, err := client.Get(ctx, shortenedId).Result()

	return redirectUrl, err
}

func DeleteRedirectUrlFromCache(ctx context.Context, shortenedId string) error {
	client := connectToCache()

	_, err := client.Del(ctx, shortenedId).Result()

	return err
}
