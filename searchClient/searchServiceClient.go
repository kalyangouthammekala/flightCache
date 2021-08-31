package searchClient

import "awesomeProject1/models"

type tfmSearchServiceClient struct {
	Type string
	ID   int
}

func (t *tfmSearchServiceClient) Search(request *models.SearchRequest) string {
	//TODO make a call to tfm airline search connector
	return ""
}
