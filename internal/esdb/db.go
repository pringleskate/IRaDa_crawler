package esdb

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
	"github.com/pkg/errors"
	"github.com/pringleskate/IRaDa_crawler/internal/models"
	"log"
	"reflect"
	"strings"
)

func (db *Database) PutIndex() error {
	ctx := context.Background()

	exists, err := db.client.IndexExists(db.index).Do(ctx)
	if err != nil {
		return errors.Wrapf(err, "Index exists failed: %s", elasticError(err))
	}

	if exists {
		return nil
	}

	_, err = db.client.CreateIndex(db.index).BodyString(mapping).Do(ctx)
	_, err = db.client.CreateIndex(db.index).Do(ctx)
	if err != nil {
		return errors.Wrapf(err, "Create index failed: %s", elasticError(err))
	}

	/*_, err = db.client.PutMapping().Index(db.index).BodyString(PutMapping).Do(ctx)
	if err != nil {
		return errors.Wrapf(err,"Put mapping failed: %s", elasticError(err))
	}*/

	return nil
}

// Create a new index.
const mapping = `
    "settings": {
        "analysis": {
            "filter": {
                "delimiter": {
                    "type": "word_delimiter",
                    "preserve_original": "true"
                },
                "jmorphy2_russian": {
                    "type": "jmorphy2_stemmer",
                    "name": "ru"
                }
            },
            "analyzer": {
                "text_ru": {
                    "tokenizer": "standard",
                    "filter": [
                        "lowercase",
                        "delimiter",
                        "jmorphy2_russian"
                    ]
                }
            }
        }
    },
    "mappings": {
            "dynamic": "strict",
            "properties":{
			"title":{
				"type": "text",
				"analyzer":"text_ru",
				"fields": {
          			"keyword": { 
            			"type": "keyword"
					}
        		}
			},
			"link":{
				"type":"keyword"
			},
			"description":{
				"type":"text",
				"fields": {
          			"keyword": { 
            			"type": "keyword"
					}
        	},
				"analyzer":"text_ru"
			},
			"pubdate":{
				"type":"date"
			},
			"article":{
				"type":"text",
				"analyzer":"text_ru"
			}
	 }
}`

func (db *Database) GetInfo() error {
	ctx := context.Background()

	info, code, err := db.client.Ping(pingAddr).Do(ctx)
	if err != nil {
		return errors.Wrapf(err,"Ping failed: %s", elasticError(err))
	}

	log.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)
	return nil
}

func (db *Database) FillDatabaseWithItems(items ...models.RepoItem) error {
	ctx := context.Background()

	for _, item := range items {
		_, err := db.client.Index().Index(db.index).Id(item.Link).BodyJson(item).Do(ctx)
		if err != nil {
			return errors.Wrapf(err, "Create index failed: %s", elasticError(err))

		}
	}

	return nil
}

func (db *Database) FindDuplicates(input ...models.Item) ([]models.Item, error) {
	ctx := context.Background()

	query := elastic.NewMatchAllQuery()
	res, err := db.client.Search().Index(db.index).Query(query).Size(lastDocsCount).Do(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "Get last 100 docs failed: %s", elasticError(err))
	}

	lastItems := map[string]models.RepoItem{}
	for _, item := range res.Each(reflect.TypeOf(models.RepoItem{})) {
		if t, ok := item.(models.RepoItem); ok {
			lastItems[t.Link] = t
		}
	}

	notDownloaded := make([]models.Item, 0)
	for _, item := range input {
		if _, ok := lastItems[item.Link]; !ok {
			notDownloaded = append(notDownloaded, item)
		}
	}

	return notDownloaded, nil
}

func (db *Database) SearchByKey(key string) error {
	ctx := context.Background()

	queryString := elastic.NewQueryStringQuery(key)
	res, err := db.client.Search().Index(db.index).Query(queryString).Do(ctx)
	if err != nil {
		return errors.Wrapf(err, "Query string for key failed: %s", elasticError(err))
	}

	if res.Hits.TotalHits.Value == 0 {
		fmt.Printf("Found no matches by key %s\n", key)
		return nil
	}

	fmt.Printf("\nFound a total of %d items by key word %s\n", res.Hits.TotalHits.Value, key)

	for _, hit := range res.Hits.Hits {
		var t models.RepoItem
		err := json.Unmarshal(hit.Source, &t)
		if err != nil {
			return errors.Wrapf(err, "Can't unmarshal response item from elastic")
		}

		fmt.Printf("Title = %s, link = %s\n", t.Title, t.Link)
	}

	return nil
}

func (db *Database) TermAggregationByField(field string) error {
	ctx := context.Background()

	aggr := elastic.NewTermsAggregation().Field(AggregationField(field)).Size(aggrSize).OrderByCountDesc()

	result, err := db.client.Search().Index(db.index).Aggregation(termAggr, aggr).Do(ctx)
	if err != nil {
		return errors.Wrapf(err, "Term aggregation failed: %s", elasticError(err))
	}

	rawMsg := result.Aggregations[termAggr]
	var ar elastic.AggregationBucketKeyItems
	err = json.Unmarshal(rawMsg, &ar)
	if err != nil {
		return errors.Wrapf(err, "Can't unmarshal aggregation response from elastic")
	}

	fmt.Printf("\nTerm aggregation results by field %s\n", field)
	for _, item := range ar.Buckets {
		fmt.Printf("%v: %v\n", item.Key, item.DocCount)
	}
	return nil
}

//Cardinality - приблизительное число уникальных значений
func (db *Database) CardinalityAggregationByField(field string) error {
	ctx := context.Background()

	aggr := elastic.NewCardinalityAggregation().Field(AggregationField(field))

	result, err := db.client.Search().Index(db.index).Aggregation(cardAggr, aggr).Do(ctx)
	if err != nil {
		return errors.Wrapf(err, "Cardinality aggregation failed: %s", elasticError(err))
	}

	rawMsg := result.Aggregations[cardAggr]
	var ar elastic.AggregationValueMetric
	err = json.Unmarshal(rawMsg, &ar)
	if err != nil {
		return errors.Wrapf(err, "Can't unmarshal aggregation response from elastic")
	}

	fmt.Printf("\nUnique values of field %s is %.0f\n", field, *ar.Value)
	return nil
}

// есть calendar interval - месяцы, недели, дни и тд, есть fixed time interval - секунды/минуты/часы/дни
func (db *Database) DateHistogramAggregation() error {
	ctx := context.Background()

	dailyAggregation := elastic.NewDateHistogramAggregation().
		Field("pubdate").
		CalendarInterval("1d").
		Format("dd.MM.YYYY")

	result, err := db.client.Search().Index(db.index).Aggregation(histAggr, dailyAggregation).Do(ctx)
	if err != nil {
		return errors.Wrapf(err, "Date histogram aggregation failed: %s", elasticError(err))
	}

	hist, found := result.Aggregations.Histogram(histAggr)
	if !found {
		return errors.Wrapf(err, "Can't find histogram with name: %s\n", histAggr)
	}
	fmt.Printf("\nIn total in date histogram there are %d buckets", len(hist.Buckets))

	for _, bucket := range hist.Buckets {
		fmt.Printf("For date %s, doccount is %d\n", *bucket.KeyAsString, bucket.DocCount)
	}

	return nil
}

func (db *Database) GetSomeItems() ([]models.RepoItem, error) {
	ctx := context.Background()

	query := elastic.NewMatchAllQuery()
	res, err := db.client.Search().Index(db.index).Query(query).Size(2).Do(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "Get last 100 docs failed: %s", elasticError(err))
	}

	items := make([]models.RepoItem, 0)
	for _, item := range res.Each(reflect.TypeOf(models.RepoItem{})) {
		if t, ok := item.(models.RepoItem); ok {
			items = append(items, t)
		}
	}

	return items, nil
}

func (db *Database) AnalyzeText(items []models.RepoItem) ([]string, error) {
	ctx := context.Background()

	analyzedItems := make([]string, 0)

	query := elastic.NewIndicesAnalyzeService(db.client).Analyzer("russian")
	for _, item := range items {
		analyzeRes, err := query.Text(item.Article).Do(ctx)
		if err != nil {
			return nil, errors.Wrapf(err, "Analyze data failed: %s", elasticError(err))
		}

		tokens := make([]string, 0)
		for _, token := range analyzeRes.Tokens {
			tokens = append(tokens, token.Token)
		}

		analyzedItems = append(analyzedItems, strings.Join(tokens, " "))
	}

	return analyzedItems, nil
}
