package store

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/asad1123/url-shortener/src/config"
	model_url "github.com/asad1123/url-shortener/src/models/url"
	"github.com/asad1123/url-shortener/src/repository/cache"
	"github.com/asad1123/url-shortener/src/repository/db"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
)

type Store struct {
	config *config.AppConfig
	db     *db.Database
	cache  *cache.Cache
}

func NewStore(config *config.AppConfig) *Store {
	db := db.NewDatabase(config.DbHostname, config.DbPort, config.DbName)
	cache := cache.NewClient(config.CacheHostname, config.CachePort, config.CachePassword, config.CacheDb)

	return &Store{config: config, db: db, cache: cache}
}

func (s *Store) SaveUrl(url model_url.Url) error {
	err := s.db.SaveNewUrl(url)
	if err != nil {
		return err
	}

	// write through cache design
	ctx := context.Background()
	err = s.cache.SaveShortenedUrlToCache(ctx, url.RedirectUrl, url.ShortenedId)

	return err
}

func (s *Store) GetUrl(id string) (string, error) {
	ctx := context.Background()
	redirectUrl, err := s.cache.GetRedirectUrlFromCache(ctx, id)
	if err != nil {
		// record our cache miss here
		msg := fmt.Sprintf("Cache miss: %s", id)
		log.Println(msg)

		url, err := s.db.GetUrl(id)

		// save this back to cache for future hits
		s.cache.SaveShortenedUrlToCache(ctx, url.RedirectUrl, url.ShortenedId)
		return url.RedirectUrl, err
	}

	return redirectUrl, err
}

func (s *Store) DeleteUrl(id string) (*mgo.ChangeInfo, error) {
	ctx := context.Background()
	err := s.cache.DeleteRedirectUrlFromCache(ctx, id)
	info, dbErr := s.db.DeleteUrl(id)

	if err != nil {
		err = errors.Wrap(err, dbErr.Error())
	} else {
		err = dbErr
	}

	return info, err
}

func (s *Store) SaveUrlUsage(usage model_url.UrlUsage) error {
	return s.db.SaveUrlUsage(usage)
}

func (s *Store) GetUrlUsage(id string, initialTimestamp time.Time) (int, error) {
	return s.db.SearchUrlUsage(id, initialTimestamp)
}
