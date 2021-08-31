package service

import (
	"awesomeProject1/models"
	"awesomeProject1/redis"
	"awesomeProject1/ruleEngine"
	"awesomeProject1/searchService"
	"encoding/json"
	"fmt"
)

type cacheService interface {
	LoadCache(knowledgeBaseDetails *models.KnowledgeBaseForCacheRule) bool
	Search(knowledgeBaseDetails *models.KnowledgeBaseForCacheRule) *models.SearchResponse
}

type FlightCacheService struct {
	Request       *models.SearchRequest
	Response      *models.SearchResponse
	SearchService searchService.SearchService
}

func (f FlightCacheService) LoadCache(knowledgeBaseDetails *models.KnowledgeBaseForCacheRule) bool {
	//load rules
	rulesLoaded := false
	kb := ruleEngine.LoadRules(knowledgeBaseDetails.Name,
		knowledgeBaseDetails.Version)
	if kb != nil && kb.RuleEntries != nil {
		//for each rule execute search and load results into the cache
		ruleEntriesMap := kb.RuleEntries
		for k, v := range ruleEntriesMap {
			fmt.Println("key :", k, " v:", v.RuleName)
		}
		rulesLoaded = true
	}

	return rulesLoaded
}

func (f FlightCacheService) Search(knowledgeBaseDetails *models.KnowledgeBaseForCacheRule) *models.SearchResponse {
	response := ruleEngine.Execute(f.Request, f.Response,
		knowledgeBaseDetails.Name, knowledgeBaseDetails.Version)
	if response.FromCache {
		//process the request by querying the cache
		cacheEntryKey := deriveCacheKeyFromRequest(f.Request)
		cacheEntry := redis.Query(cacheEntryKey)
		fmt.Println("Cache entry retrieved: ", cacheEntry.Value)
		if cacheEntry.Value == "" {
			//call search service
			f.requestSearchService(response, cacheEntryKey, true)

		} else {
			response.FromCache = true
			loadResponseWithResult(cacheEntry.Value, "", response, false)
			fmt.Println("Result loaded from cache with entry: ", cacheEntryKey)
		}

	} else {
		//get from search service
		response = f.requestSearchService(response, "", false)
	}
	return response
}

func loadResponseWithResult(result, cacheEntryKey string, response *models.SearchResponse, addToCache bool) *models.SearchResponse {
	tfmRespose := &models.TfmResponse{}
	error := json.Unmarshal([]byte(result), &tfmRespose)
	if error != nil {
		fmt.Println("Unable to unmashal string response to tfmresponse")
	}
	response.TfmRessponse = *tfmRespose

	//check and load into cache
	//store the result in cache for future use
	if result != "" && addToCache && cacheEntryKey != "" {
		redis.AddEntry(cacheEntryKey, result)
		fmt.Println("Cache entry added: ", cacheEntryKey)
	}

	return response
}

func (f FlightCacheService) requestSearchService(response *models.SearchResponse, cacheEntryKey string, addToCache bool) *models.SearchResponse {
	response.FromCache = false
	result := f.SearchService.Search(f.Request)
	loadResponseWithResult(result, cacheEntryKey, response, addToCache)
	return response
}

func deriveCacheKeyFromRequest(request *models.SearchRequest) string {
	key := ""
	journeyType := ""
	//JFKHAJ-R-20122021-26122021
	if request.RoundTrip {
		journeyType = "R"
	} else {
		journeyType = "O"
	}
	key = request.DepartureAirportCode + request.ArrivalAirportCode + "-" + journeyType + "" +
		"-" + models.TransformDate(request.DepartureDateTime) + "-" +
		models.TransformDate(request.ArrivalDateTime)

	fmt.Println("Key for Cache ", key)

	return key
}
