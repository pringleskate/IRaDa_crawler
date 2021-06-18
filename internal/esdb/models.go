package esdb

import (
	"github.com/olivere/elastic/v7"
	"log"
)

type Database struct {
	client *elastic.Client
	index string
}

func NewDatabase(cl *elastic.Client) Database {
	return Database{client: cl, index: indexName}
}

func NewClient() (*elastic.Client, error){
	client, err := elastic.NewSimpleClient()
	if err != nil {
		return nil, err
	}

	return client, nil
}

const (
	pingAddr = "http://127.0.0.1:9200"

	indexName = "items"
	termAggr  = "term_aggr"
	histAggr  = "date_hist"
	cardAggr  = "card_aggr"

	lastDocsCount = 50
	aggrSize = 20
)

func AggregationField(field string) string {
	if field == "title" || field == "description" {
		return field + ".keyword"
	}

	return field
}

func elasticError(err error) string {
	e, ok := err.(*elastic.Error)
	if !ok {
		log.Println("Can't convert to *elastic.Error")
		return ""
	}

	log.Printf("Elastic failed with status %d and error %s", e.Status, e.Details.Reason)
	return e.Details.Reason
}

const PutMapping = `
	{
	  "properties":{
			"title":{
				"type": "text",
				"analyzer":"russian",
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
				"analyzer":"russian"
			},
			"pubdate":{
				"type":"date"
			},
			"article":{
				"type":"text",
				"analyzer":"russian"
			}
	 }
}`

const GetDocuments = `
	"query": {
		  "match_all": {}
	   }`
