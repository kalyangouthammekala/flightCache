package service

import (
	"awesomeProject1/models"
	"awesomeProject1/redis"
	"awesomeProject1/ruleEngine"
	"fmt"
)

type cacheService interface {
	LoadCache(knowledgeBaseDetails *models.KnowledgeBaseForCacheRule) bool
	Search(knowledgeBaseDetails *models.KnowledgeBaseForCacheRule) *models.SearchResponse
}

type FlightCacheService struct {
	Request  *models.SearchRequest
	Response *models.SearchResponse
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
		cacheEntry := redis.Query(deriveCacheKeyFromRequest(f.Request))
		fmt.Println("Cache entry retrieved: ", cacheEntry.Value)
	} else {

	}
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
