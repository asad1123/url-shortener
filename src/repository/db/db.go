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

type Database struct {
	session *mgo.Database
}

func connectToDb(host string, port string, dbName string) (*mgo.Database, error) {
	url := fmt.Sprintf("%s:%s", host, port)

	session, err := mgo.Dial(url)
	if err != nil {
		return nil, err
	}

	database := session.DB(dbName)
	return database, nil
}

func NewDatabase(host string, port string, dbName string) *Database {
	db, err := connectToDb(host, port, dbName)
	if err != nil {
		fmt.Println(err)
	}

	return &Database{session: db}
}

func (d *Database) SaveNewUrl(url model_url.Url) error {
	db := d.session

	err := db.C(UrlCollection).Insert(url)
	return err
}

func (d *Database) GetUrl(id string) (model_url.Url, error) {
	db := d.session

	url := model_url.Url{}
	err := db.C(UrlCollection).Find(bson.M{"shortenedid": id}).One(&url)

	return url, err
}

func (d *Database) DeleteUrl(id string) (*mgo.ChangeInfo, error) {
	db := d.session

	info, err := db.C(UrlCollection).RemoveAll(bson.M{"shortenedid": id})

	return info, err
}

func (d *Database) SaveUrlUsage(usage model_url.UrlUsage) error {
	db := d.session

	err := db.C(AnalyticsCollection).Insert(usage)
	return err
}

func (d *Database) SearchUrlUsage(id string, since time.Time) (int, error) {
	db := d.session
	count, err := db.C(AnalyticsCollection).Find(bson.M{"shortenedid": id, "accessedat": bson.M{"$gte": since}}).Count()

	return count, err
}
