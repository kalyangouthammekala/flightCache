package test

import (
	"awesomeProject1/mocks"
	"awesomeProject1/models"
	"awesomeProject1/redis"
	"awesomeProject1/searchService"
	"awesomeProject1/service"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"path/filepath"
	"testing"
	"time"
)

var ctrl *gomock.Controller

func getMockSearchService(t *testing.T, request *models.SearchRequest, addCacheEntry bool) *mocks.MockSearchService {
	ctrl = gomock.NewController(t)

	//defer ctrl.Finish()
	mockSearchService := mocks.NewMockSearchService(ctrl)
	sampleSearchResultFile, err := filepath.Abs("sampleSearchResult.json")
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

	if addCacheEntry {
		redis.AddEntry(models.DeriveCacheKeyFromRequest(request), string(sampleSearchResultJson), nil)
	}

	return mockSearchService
}

func TestSearchWithEmptyCache(t *testing.T) {
	request := &models.SearchRequest{
		AirlineCode:          "KL",
		DepartureAirportCode: "AMS",
		ArrivalAirportCode:   "NYC",
		DepartureDateTime: time.Date(
			2021,
			11,
			19,
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
	}
	mockSearchService := getMockSearchService(t, request, false)
	defer ctrl.Finish()
	//remove all entries in cache
	//redis.RemoveAllEntries()

	//cache service
	flightCacheService := &service.FlightCacheService{
		Request:       request,
		Response:      &models.SearchResponse{},
		SearchService: mockSearchService,
	}

	kbDetails := &models.KnowledgeBaseForCacheRule{
		Name:    "Test",
		Version: "0.0.1",
	}
	response := flightCacheService.Search(kbDetails, nil)

	assert.Equal(t, false, response.FromCache, "Cache should not contain it!!")
}

func TestSearchWithCacheEntry(t *testing.T) {
	request := &models.SearchRequest{
		AirlineCode:          "KL",
		DepartureAirportCode: "AMS",
		ArrivalAirportCode:   "NYC",
		DepartureDateTime: time.Date(
			2021,
			11,
			19,
			12,
			35,
			0,
			0,
			time.UTC,
		),
		ArrivalDateTime: time.Date(
			2021,
			9,
			06,
			12,
			35,
			0,
			0,
			time.UTC,
		),
		RoundTrip:   true,
		BookingTime: time.Now(),
	}

	//cache service
	searchService := &searchService.DummySearchServiceImpl{}
	flightCacheService := &service.FlightCacheService{
		Request: request,
		Response: &models.SearchResponse{
			FromCache:            false,
			AirlineCode:          request.AirlineCode,
			DepartureAirportCode: request.DepartureAirportCode,
			ArrivalAirportCode:   request.ArrivalAirportCode,
			RoundTrip:            request.RoundTrip,
			BookingTime:          request.BookingTime,
		},
		SearchService: searchService,
	}

	kbDetails := &models.KnowledgeBaseForCacheRule{
		Name:    "Test",
		Version: "0.0.1",
	}
	//add cache entry
	key := models.DeriveCacheKeyFromRequest(request)
	redis.AddEntry(key, searchService.Search(request), nil)
	response := flightCacheService.Search(kbDetails, nil)

	assert.Equal(t, true, response.FromCache, "Cache should contain the key ", key)
}
