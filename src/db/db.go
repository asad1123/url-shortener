package db

import (
	"fmt"
	"time"

	model_url "github.com/asad1123/url-shortener/src/models/url"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const UrlCollection = "url"
const AnalyticsCollection = "url-analytics"

func connectToDb() (*mgo.Database, error) {
	host := "localhost"
	databaseName := "urls"

	session, err := mgo.Dial(host)
	if err != nil {
		return nil, err
	}

	database := session.DB(databaseName)
	return database, nil
}

func getDbHandle() *mgo.Database {
	db, err := connectToDb()
	if err != nil {
		fmt.Println(err)
	}

	return db
}

func SaveNewUrl(url model_url.Url) error {
	db := *getDbHandle()

	err := db.C(UrlCollection).Insert(url)
	return err
}

func GetUrl(id string) (model_url.Url, error) {
	db := *getDbHandle()

	url := model_url.Url{}
	err := db.C(UrlCollection).Find(bson.M{"shortenedid": id}).One(&url)

	return url, err
}

func DeleteUrl(id string) (*mgo.ChangeInfo, error) {
	db := *getDbHandle()

	info, err := db.C(UrlCollection).RemoveAll(bson.M{"shortenedid": id})

	return info, err
}

func SaveUrlUsage(usage model_url.UrlUsage) error {
	db := *getDbHandle()

	err := db.C(AnalyticsCollection).Insert(usage)
	return err
}

func SearchUrlUsage(id string, since time.Time) (int, error) {
	db := *getDbHandle()

	count, err := db.C(AnalyticsCollection).Find(bson.M{"shortenedid": id, "accessedat": bson.M{"$gte": since}}).Count()

	return count, err
}
