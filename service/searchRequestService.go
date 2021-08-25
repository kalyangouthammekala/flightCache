package service

import "awesomeProject1/models"

//Service to make a request to get search results from external service
type searchService interface {
	search(request *models.SearchRequest) string
}

type tfmSearchServiceClient struct {
}

func (tfmSearchServiceClient) search(request *models.SearchRequest) {
	//TODO make a call to tfm airline search connector
}
