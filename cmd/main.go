package main

import (
	"fmt"
	"github.com/pringleskate/IRaDa_crawler/internal/esdb"
	"github.com/pringleskate/IRaDa_crawler/internal/minhash"
	"github.com/pringleskate/IRaDa_crawler/internal/models"
	"github.com/pringleskate/IRaDa_crawler/internal/parser"
	"log"
	"sync"
)

func main() {
	// парсинг RSS
	items, err := parser.ParseRSS()
	if err != nil {
		log.Print(err)
		return
	}

	// создаем клиент для подключения к базе
	cl, err := esdb.NewClient()
	if err != nil {
		log.Println(err)
		return
	}

	db := esdb.NewDatabase(cl)


	err = db.PutIndex()
	if err != nil {
		log.Println(err)
		return
	}

	notDownloaded, err := db.FindDuplicates(items...)
	if err != nil {
		log.Println(err)
		return
	}

	repoItems := make([]models.RepoItem, 0)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, item := range notDownloaded {
		wg.Add(1)

		it := item

		go func() {
			defer wg.Done()

			article, err := parser.ParseArticle(it.Link)
			if err != nil {
				return
			}

			mu.Lock()
			repoItems = append(repoItems, models.FillRepoItem(it, article))
			mu.Unlock()
		}()


		wg.Wait()
	}

	for _, art := range repoItems {
		fmt.Println("Downloaded one more article: ", art.Link)
	}

	err = db.FillDatabaseWithItems(repoItems...)
	if err != nil {
		log.Println(err)
		return
	}

	//-----------------

	err = db.SearchByKey("бесы")
	if err != nil {
		log.Println(err)
		return
	}

	err = db.TermAggregationByField("title")
	if err != nil {
		log.Println(err)
		return
	}

	err = db.CardinalityAggregationByField("title")
	if err != nil {
		log.Println(err)
		return
	}

	err = db.CardinalityAggregationByField("pubdate")
	if err != nil {
		log.Println(err)
		return
	}

	err = db.DateHistogramAggregation()
	if err != nil {
		log.Println(err)
		return
	}



	//-----------------
	itms, err := db.GetSomeItems()
	if err != nil {
		log.Println(err)
		return
	}

	analyzedItems, err := db.AnalyzeText(itms)
	if err != nil {
		log.Println(err)
		return
	}

	first := minhash.SplitShingles(analyzedItems[0])
	second := minhash.SplitShingles(analyzedItems[1])

	coef := minhash.CompareByMinHash(first, second)
	fmt.Println("Coef Jak: ", coef)
}
