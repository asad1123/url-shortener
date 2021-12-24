package db

import (
	"fmt"

	model_url "github.com/asad1123/url-shortener/src/models/url"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const UrlCollection = "url"

func ConnectToDb() (*mgo.Database, error) {
	host := "localhost"
	databaseName := "urls"

	session, err := mgo.Dial(host)
	if err != nil {
		return nil, err
	}

	database := session.DB(databaseName)
	return database, nil
}

func GetDbHandle() *mgo.Database {
	db, err := ConnectToDb()
	if err != nil {
		fmt.Println(err)
	}

	return db
}

func SaveNewUrl(url model_url.Url) error {
	db := *GetDbHandle()

	err := db.C(UrlCollection).Insert(url)
	return err
}

func GetUrl(id string) (model_url.Url, error) {
	db := *GetDbHandle()

	url := model_url.Url{}
	err := db.C(UrlCollection).Find(bson.M{"shortenedid": id}).One(&url)

	return url, err
}
