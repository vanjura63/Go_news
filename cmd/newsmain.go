package main

import (
	"Go_news/pkg/api"
	"Go_news/pkg/dbnews"
	"Go_news/pkg/rss"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Application configuration
type config struct {
	URLS   []string `json:"rss"`
	Period int      `json:"request_period"`
}

func main() {

	connstr := "postgres://postgres:ts950sdx@localhost/dbnews"

	// Init DB

	db, err := dbnews.New(connstr)
	if err != nil {
		log.Fatal(err)
	}
	api := api.New(db)

	// Read fiel configuration

	b, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatal(err)
	}
	var config config
	err = json.Unmarshal(b, &config)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Файл конфига:", b)
	//fmt.Println("DB:", db)

	//Parsing news in a separate stream

	chPosts := make(chan []dbnews.Post)
	chErrs := make(chan error)
	for _, url := range config.URLS {
		go parseURL(url, db, chPosts, chErrs, config.Period)
	}

	// Recording a news stream in DB
	go func() {
		for posts := range chPosts {
			db.StoreNews(posts)
		}
	}()

	// Control error
	go func() {
		for err := range chErrs {
			log.Println("Errors:", err)
		}
	}()

	// Web server start
	err = http.ListenAndServe(":80", api.Router())
	if err != nil {
		log.Fatal(err)
	}

}

// Flow chanals
func parseURL(url string, db *dbnews.DB, posts chan<- []dbnews.Post, errs chan<- error, period int) {
	for {
		news, err := rss.Parse(url)
		if err != nil {
			errs <- err
			continue
		}
		posts <- news
		time.Sleep(time.Minute * time.Duration(period))
	}
}
