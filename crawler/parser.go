package crawler

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func Parse() {
	resp, err := http.Get("https://lenta.ru/rss/articles")
	if resp != nil && resp.StatusCode != http.StatusOK {
		log.Println(resp.Status)
		return
	}

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

	var rssData RSS

	err = xml.Unmarshal(body, &rssData)
	if err != nil {
		log.Println(err)
		return
	}

	PrintResult(rssData.Channel.Item...)
}

func PrintResult(items ...Item) {
	for _, item := range items {
		fmt.Println("GUID:", item.Guid)
		fmt.Println("Title:", item.Title)
		fmt.Println("Link:", item.Link)
		fmt.Println("Description:", item.Description)
		fmt.Println("Pubdate:", item.Pubdate)
		fmt.Println("Category:", item.Category)
	}
}