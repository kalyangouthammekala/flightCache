package service

import (
	"awesomeProject1/models"
	"awesomeProject1/redis"
	"awesomeProject1/searchService"
	"encoding/json"
	"fmt"
	"github.com/magiconair/properties"
	"log"
)

type cacheService interface {
	LoadCache(knowledgeBaseDetails *models.KnowledgeBaseForCacheRule) bool
	Search(knowledgeBaseDetails *models.KnowledgeBaseForCacheRule) *models.SearchResponse
}

type FlightCacheService struct {
	Request  *models.SearchRequest
	Response *models.SearchResult
	//Response	*models.SearchResult
	SearchService searchService.SearchService
}

//func (f FlightCacheService) LoadCache(knowledgeBaseDetails *models.KnowledgeBaseForCacheRule) bool {
//	//load rules
//	rulesLoaded := false
//	kb := ruleEngine.LoadRules(knowledgeBaseDetails.Name,
//		knowledgeBaseDetails.Version)
//	if kb != nil && kb.RuleEntries != nil {
//		//for each rule execute search and load results into the cache
//		ruleEntriesMap := kb.RuleEntries
//		for k, v := range ruleEntriesMap {
//			fmt.Println("key :", k, " v:", v.RuleName)
//		}
//		rulesLoaded = true
//	}
//
//	return rulesLoaded
//}

//func (f FlightCacheService) Search(knowledgeBaseDetails *models.KnowledgeBaseForCacheRule, flightCacheProperties *properties.Properties) *models.SearchResponse {
//	//response := ruleEngine.Execute(f.Request, f.Response,
//	//	knowledgeBaseDetails.Name, knowledgeBaseDetails.Version)
//
//	response := ruleEngine.RuleEngineSearchResponse(f.Request, f.Response)
//	if response.FromCache {
//		//process the request by querying the cache
//		cacheEntryKey := models.DeriveCacheKeyFromRequest(f.Request)
//		cacheEntry := redis.Query(cacheEntryKey, flightCacheProperties)
//		fmt.Println("Cache entry retrieved: ", cacheEntry.Value)
//		if cacheEntry.Value == "" {
//			//call search service
//			f.RequestSearchService(response, cacheEntryKey, true, flightCacheProperties)
//
//		} else {
//			response.FromCache = true
//			LoadResponseWithResult(cacheEntry.Value, "", response, false, flightCacheProperties)
//			fmt.Println("Result loaded from cache with entry: ", cacheEntryKey)
//		}
//
//	} else {
//		//get from search service
//		response = f.RequestSearchService(response, "", false, flightCacheProperties)
//	}
//	return response
//}

func LoadResultFromCache(cacheValue string) (*models.Result, error) {
	var result *models.Result
	err := json.Unmarshal([]byte(cacheValue), &result)
	if err != nil {
		log.Println("Unable to unmashal string response to tfmresponse")
	}

	return result, err

}

func LoadResponseWithResult(result, cacheEntryKey string, response *models.SearchResult, addToCache bool, p *properties.Properties) *models.SearchResult {
	tfmResponse := &models.Result{}
	err := json.Unmarshal([]byte(result), &tfmResponse)
	if err != nil {
		fmt.Println("Unable to unmashal string response to tfmresponse")
	}
	response.Result = *tfmResponse

	//check and load into cache
	//store the result in cache for future use
	if result != "" && addToCache && cacheEntryKey != "" {
		redis.AddEntry(cacheEntryKey, result, p)
		fmt.Println("Cache entry added: ", cacheEntryKey)
	}

	return response
}

func (f FlightCacheService) RequestSearchService(response *models.SearchResult, cacheEntryKey string, addToCache bool, p *properties.Properties) *models.SearchResult {
	response.AdditionalInfo.NeedToBeCached = false
	result := f.SearchService.Search(f.Request)
	LoadResponseWithResult(result, cacheEntryKey, response, addToCache, p)
	return response
}
