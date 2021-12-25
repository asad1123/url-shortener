package store

import (
	"context"
	"fmt"
	"log"
	"time"

	model_url "github.com/asad1123/url-shortener/src/models/url"
	"github.com/asad1123/url-shortener/src/repository/cache"
	"github.com/asad1123/url-shortener/src/repository/db"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
)

func SaveUrl(url model_url.Url) error {
	err := db.SaveNewUrl(url)
	if err != nil {
		return err
	}

	// write through cache design
	ctx := context.Background()
	err = cache.SaveShortenedUrlToCache(ctx, url.RedirectUrl, url.ShortenedId)

	return err
}

func GetUrl(id string) (string, error) {
	ctx := context.Background()
	redirectUrl, err := cache.GetRedirectUrlFromCache(ctx, id)
	if err != nil {
		// record our cache miss here
		msg := fmt.Sprintf("Cache miss: %s", id)
		log.Println(msg)

		url, err := db.GetUrl(id)

		// save this back to cache for future hits
		cache.SaveShortenedUrlToCache(ctx, url.RedirectUrl, url.ShortenedId)
		return url.RedirectUrl, err
	}

	return redirectUrl, err
}

func DeleteUrl(id string) (*mgo.ChangeInfo, error) {
	ctx := context.Background()
	err := cache.DeleteRedirectUrlFromCache(ctx, id)
	info, dbErr := db.DeleteUrl(id)

	if err != nil {
		err = errors.Wrap(err, dbErr.Error())
	} else {
		err = dbErr
	}

	return info, err
}

func SaveUrlUsage(usage model_url.UrlUsage) error {
	return db.SaveUrlUsage(usage)
}

func GetUrlUsage(id string, initialTimestamp time.Time) (int, error) {
	return db.SearchUrlUsage(id, initialTimestamp)
}
