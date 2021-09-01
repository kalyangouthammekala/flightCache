package searchService

import (
	"awesomeProject1/models"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

// SearchService Service to make a request to get search results from external service
type SearchService interface {
	Search(request *models.SearchRequest) string
}

type DummySearchServiceImpl struct {
}

func (d *DummySearchServiceImpl) Search(request *models.SearchRequest) string {
	sampleSearchResultFile, err := filepath.Abs("../resources/sampleSearchResult.json")
	if err != nil {
		fmt.Println("Couldn't parse sample search result file path")
		panic(err)
	}
	sampleSearchResultJson, err := ioutil.ReadFile(sampleSearchResultFile)
	if err != nil {
		fmt.Println("Couldn't read sample search result")
		panic(err)
	}
	return string(sampleSearchResultJson)
}
