package searchService

import (
	"awesomeProject1/models"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// SearchService Service to make a request to get search results from external service
type SearchService interface {
	Search(request *models.SearchRequest) string
}

type DummySearchServiceImpl struct {
}

func (d *DummySearchServiceImpl) Search(request *models.SearchRequest) string {
	sampleSearchResultFile, err := filepath.Abs("sampleSearchResult.json")
	if err != nil {
		fmt.Println("Couldn't parse sample search result file path")
		panic(err)
	}
	sampleSearchResultJson, err := ioutil.ReadFile(sampleSearchResultFile)
	if err != nil {
		fmt.Println("Couldn't read sample search result")
		if strings.Contains(err.Error(), "The system cannot find the path specified") {
			err = nil
			sampleSearchResultFile, err = filepath.Abs("sampleSearchResult.json")
			sampleSearchResultJson, err = ioutil.ReadFile(sampleSearchResultFile)
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
	return string(sampleSearchResultJson)
}

type tfmSearchServiceClient struct {
}

func (t *tfmSearchServiceClient) Search(request *models.SearchRequest) string {
	//TODO make a call to tfm airline search connector
	return ""
}
