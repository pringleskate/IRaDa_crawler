package models

import (
	"encoding/xml"
	"time"
)

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
	Version string   `xml:"version,attr"`
}

type Channel struct {
	XMLName     xml.Name `xml:"channel"`
	Item        []Item   `xml:"item"`
}

type Item struct {
	XMLName     xml.Name  `xml:"item"`
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Pubdate     string    `xml:"pubDate"`
}

type RepoItem struct {
	Title       string    `json:"title"` //type text
	Link        string    `json:"link"` //type keyword
	Description string    `json:"description"` //type text
	Pubdate     time.Time  `json:"pubdate"` //type date
	Article     string    `json:"article"`  //type text
}

func FillRepoItem(item Item, article string) RepoItem {
	date, _ := time.Parse(time.RFC1123Z, item.Pubdate)
	return RepoItem{
		Title:       item.Title,
		Link:        item.Link,
		Description: item.Description,
		Pubdate:     date,
		Article:     article,
	}
}
