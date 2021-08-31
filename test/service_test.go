package test

import (
	"awesomeProject1/mocks"
	"awesomeProject1/models"
	"awesomeProject1/redis"
	"awesomeProject1/service"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"path/filepath"
	"testing"
	"time"
)

func getMockSearchService(t *testing.T) *mocks.MockSearchService {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()
	mockSearchService := mocks.NewMockSearchService(ctrl)
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
	mockSearchService.EXPECT().Search(gomock.Any()).Return(string(sampleSearchResultJson))

	return mockSearchService
}

func TestFlightCacheService_Search(t *testing.T) {

	mockSearchService := getMockSearchService(t)
	//remove all entries in cache
	redis.RemoveAllEntries()

	//cache service
	flightCacheService := &service.FlightCacheService{
		Request: &models.SearchRequest{
			AirlineCode:          "KL",
			DepartureAirportCode: "AMS",
			ArrivalAirportCode:   "NYC",
			DepartureDateTime: time.Date(
				2021,
				9,
				04,
				12,
				35,
				0,
				0,
				time.UTC,
			),
			ArrivalDateTime: time.Date(
				2021,
				9,
				05,
				12,
				35,
				0,
				0,
				time.UTC,
			),
			RoundTrip:   true,
			BookingTime: time.Now(),
		},
		Response:      &models.SearchResponse{},
		SearchService: mockSearchService,
	}

	kbDetails := &models.KnowledgeBaseForCacheRule{
		Name:    "Test",
		Version: "0.0.1",
	}
	response := flightCacheService.Search(kbDetails)

	assert.Equal(t, false, response.FromCache, "Cache should not contain it!!")
}
