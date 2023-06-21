package service

import (
	"gopkg.in/olivere/elastic.v6"
)

var (
	esClient *elastic.Client
)

func NewEsClient(hosts []string, user, password string) (*elastic.Client, error) {
	esClient, err := elastic.NewClient(
		elastic.SetURL(hosts...),
		// elastic.SetSniff(false),
		// elastic.SetHealthcheck(false),
		elastic.SetBasicAuth(user, password),
	)
	if err != nil {
		return nil, err
	}
	return esClient, nil
}
