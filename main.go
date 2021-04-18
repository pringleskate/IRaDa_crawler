package main

import (
	"crawler/crawler"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	resp, err := http.Get("https://lenta.ru/rss/articles")
	log.Println(resp.Status)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}


	var rssData crawler.RSS

	err = xml.Unmarshal(body, &rssData)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Printf("%+v\n", rssData)
}
