package service

import (
	"awesomeProject1/models"
	"awesomeProject1/redis"
	"awesomeProject1/ruleEngine"
	"awesomeProject1/searchService"
	"encoding/json"
	"fmt"
	"github.com/magiconair/properties"
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

func (f FlightCacheService) Search(knowledgeBaseDetails *models.KnowledgeBaseForCacheRule, flightCacheProperties *properties.Properties) *models.SearchResponse {
	//response := ruleEngine.Execute(f.Request, f.Response,
	//	knowledgeBaseDetails.Name, knowledgeBaseDetails.Version)

	response := ruleEngine.RuleEngineSearchResponse(f.Request, f.Response)
	if response.FromCache {
		//process the request by querying the cache
		cacheEntryKey := models.DeriveCacheKeyFromRequest(f.Request)
		cacheEntry := redis.Query(cacheEntryKey, flightCacheProperties)
		fmt.Println("Cache entry retrieved: ", cacheEntry.Value)
		if cacheEntry.Value == "" {
			//call search service
			f.requestSearchService(response, cacheEntryKey, true, flightCacheProperties)

		} else {
			response.FromCache = true
			loadResponseWithResult(cacheEntry.Value, "", response, false, flightCacheProperties)
			fmt.Println("Result loaded from cache with entry: ", cacheEntryKey)
		}

	} else {
		//get from search service
		response = f.requestSearchService(response, "", false, flightCacheProperties)
	}
	return response
}

func loadResponseWithResult(result, cacheEntryKey string, response *models.SearchResponse, addToCache bool, p *properties.Properties) *models.SearchResponse {
	tfmRespose := &models.TfmResponse{}
	err := json.Unmarshal([]byte(result), &tfmRespose)
	if err != nil {
		fmt.Println("Unable to unmashal string response to tfmresponse")
	}
	response.TfmRessponse = *tfmRespose

	//check and load into cache
	//store the result in cache for future use
	if result != "" && addToCache && cacheEntryKey != "" {
		redis.AddEntry(cacheEntryKey, result, p)
		fmt.Println("Cache entry added: ", cacheEntryKey)
	}

	return response
}

func (f FlightCacheService) requestSearchService(response *models.SearchResponse, cacheEntryKey string, addToCache bool, p *properties.Properties) *models.SearchResponse {
	response.FromCache = false
	result := f.SearchService.Search(f.Request)
	loadResponseWithResult(result, cacheEntryKey, response, addToCache, p)
	return response
}
