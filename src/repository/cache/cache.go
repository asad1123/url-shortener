package cache

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type Cache struct {
	client *redis.Client
}

func connectToCache(addr string, password string, db int) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return client
}

func NewClient(cacheHostname string, cachePort string, cachePassword string, cacheDb int) *Cache {
	addr := fmt.Sprintf("%s:%s", cacheHostname, cachePort)
	return &Cache{connectToCache(addr, cachePassword, cacheDb)}
}

func (c *Cache) SaveShortenedUrlToCache(ctx context.Context, redirectUrl string, shortenedId string) error {
	client := c.client

	_, err := client.Set(ctx, shortenedId, redirectUrl, 0).Result()

	return err
}

func (c *Cache) GetRedirectUrlFromCache(ctx context.Context, shortenedId string) (string, error) {
	client := c.client

	redirectUrl, err := client.Get(ctx, shortenedId).Result()

	return redirectUrl, err
}

func (c *Cache) DeleteRedirectUrlFromCache(ctx context.Context, shortenedId string) error {
	client := c.client

	_, err := client.Del(ctx, shortenedId).Result()

	return err
}
