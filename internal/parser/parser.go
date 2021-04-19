package parser

import (
	"encoding/xml"
	"fmt"
	"github.com/pkg/errors"
	"github.com/pringleskate/IRaDa_crawler/internal/models"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// функция парсинга rss-ленты
func ParseRSS() ([]models.Item, error) {
	// скачиваем данные страницы
	resp, err := http.Get("https://lenta.ru/rss/articles")
	if resp != nil && resp.StatusCode != http.StatusOK {
		log.Println(resp.Status)
		return nil, err
	}

	if err != nil {
		return nil, errors.Wrapf(err, "Get RSS failed")
	}
	defer resp.Body.Close()

	// читаем сам ответ (тело запроса)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "Read body failed")
	}

	var rssData models.RSS

	// переводим байты в вид структур, которыми сможем пользоваться в коде
	err = xml.Unmarshal(body, &rssData)
	if err != nil {
		return nil, errors.Wrapf(err, "Unmarshaling failed")
	}

	return rssData.Channel.Item, nil
}

func PrintResult(items ...models.Item) {
	for _, item := range items {
		fmt.Println("Title:", item.Title)
		fmt.Println("Link:", item.Link)
		fmt.Println("Description:", item.Description)
		fmt.Println("Pubdate:", item.Pubdate)
	}
}

func ParseArticle(link string) (string, error) {
	// скачиваем страницу
	res, err := http.Get(link)
	if err != nil {
		return "", errors.Wrapf(err, "GET-request by link failed")
	}
	defer res.Body.Close()

	// используем библиотеку goquery, формирует скачанные данные в формат своего документа
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", errors.Wrapf(err, "Getting doc from body failed")
	}

	article := ""
	// делаем поиск по CSS-селектору
	doc.Find(".js-topic__text").Each(func(i int, s *goquery.Selection) {
		article += s.Find("p").Text()
	})

	// возвращаем саму статью в виде сырой строки
	return article, nil
}