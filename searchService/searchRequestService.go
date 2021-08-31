package searchService

import (
	"awesomeProject1/models"
	"fmt"
)

// SearchService Service to make a request to get search results from external service
type SearchService interface {
	Search(request *models.SearchRequest) string
}

func Test(service SearchService) {
	fmt.Println("Test the interface implementation!!")
}
