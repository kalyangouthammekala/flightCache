package controller

import (
	"awesomeProject1/models"
	"awesomeProject1/redis"
	"awesomeProject1/ruleEngine"
	"awesomeProject1/service"
	"encoding/json"
	"fmt"
	"github.com/magiconair/properties"
	"log"
	"net/http"
	"strconv"
	"time"
)

var FlightCacheProperties *properties.Properties

//func StartServer() {
//	http.HandleFunc("/", handler)
//	http.HandleFunc("/flightCache/search", searchHandler)
//	log.Fatal(http.ListenAndServe(":8081", nil))
//}

func handler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "Hi there, %s!", r.URL.Path[1:])
	if err != nil {
		fmt.Println("Error in handling request ", err.Error())
	}
}

//func SearchHandler(w http.ResponseWriter, r *http.Request) {
func SearchHandler(request string, flightCacheProperties *properties.Properties) ([]byte, error) {

	log.Println("Entering into SearchHandler")
	var wideSearchRequest models.WideSearchQuery
	var response []byte
	err := json.Unmarshal([]byte(request), &wideSearchRequest)
	if err != nil {
		panic(err)
	}

	searchResponse := process(wideSearchRequest)
	response, err = json.Marshal(searchResponse)
	return response, err
}

func processFltCacheRuleConfigRequest(request *models.SearchRequest, resultMap map[*models.SearchRequest]*models.SearchResponseFromRuleEngine) {
	searchRuleEngineResponse := make(chan *models.SearchResponseFromRuleEngine)
	go ruleEngine.RuleEngineClientResponseThroughChan(request, searchRuleEngineResponse)
	response := <-searchRuleEngineResponse
	resultMap[request] = response
}

func processFltCacheSearch(request *models.SearchRequest, resultsMap map[string]models.SearchResult) {
	/*searchResult := make(chan *models.SearchResult)
	searchKey := make(chan *string)
	go performSearchThroughChan(request, FlightCacheProperties, searchResult, searchKey)
	resultKey := <-searchKey
	result := <-searchResult
	resultsMap[*resultKey] = *result*/
	result, key, _ := performSearch(request, FlightCacheProperties)
	resultsMap[*key] = *result
}

func process(request models.WideSearchQuery) models.FlightCacheServiceResponse {

	log.Println("Entering into Process")

	resultsMap := make(map[string]models.SearchResult)

	//split the request into individual requests
	subRequests := splitRequestIntoUnits(&request)

	log.Println("No of sub request", len(subRequests))

	//TODO check for batch processing
	/*batchSizePropertyVal, isBatchSizePropertyPresent := FlightCacheProperties.Get("batch-size")
	batchProcessingEnabledPropertyVal, isBatchProcPropertyPresent := FlightCacheProperties.Get("batch-proc-enabled")
	if isBatchProcPropertyPresent && isBatchSizePropertyPresent {
		isBatchProcEnabled, err := strconv.ParseBool(batchProcessingEnabledPropertyVal)
		batchSize, err := strconv.Atoi(batchSizePropertyVal)
		if err != nil {
			log.Println("Invalid property values for batch processing in configuration")
		}
		log.Println(isBatchProcEnabled, batchSize)
	}*/

	for _, subRequest := range subRequests {
		flightCacheRuleConfigurationCheckProperty, _ := FlightCacheProperties.Get("isFlightCacheRuleConfigurationCheckEnabled")

		isFlightCacheRuleConfigurationCheckEnabled, _ := strconv.ParseBool(flightCacheRuleConfigurationCheckProperty)

		log.Println("Rule Engine check Enabled: ", isFlightCacheRuleConfigurationCheckEnabled)
		//If rule configuration check is enabled
		if isFlightCacheRuleConfigurationCheckEnabled {
			resultMap := make(map[*models.SearchRequest]*models.SearchResponseFromRuleEngine)

			//Flt Cache rule engine check processing
			log.Println("Calling Rule Engine per request")
			processFltCacheRuleConfigRequest(subRequest, resultMap)

			//Flight Cache search processing
			log.Println("Calling cache search with resultMap with entries count: ", len(resultMap))
			for mapEntryRequest, ruleConfigResponse := range resultMap {
				log.Println("Cacheable : ", ruleConfigResponse.Cacheable)
				if ruleConfigResponse.Cacheable {
					processFltCacheSearch(mapEntryRequest, resultsMap)
				} else {
					log.Println("Cacheable : ", ruleConfigResponse.Cacheable)
					searchKey := models.DeriveCacheKeyFromRequest(mapEntryRequest)
					log.Println("searchKey : ", searchKey)
					result := &models.SearchResult{
						Result: models.Result{},
						AdditionalInfo: models.AdditionalInfo{
							NeedToBeCached: false,
							ResultType:     "",
						},
					}
					resultsMap[searchKey] = *result
				}
			}
			/*go func() {
				for mapEntryRequest, ruleConfigResponse := range resultMap {
					if ruleConfigResponse.Cacheable {
						go processFltCacheSearch(mapEntryRequest, resultsMap)
					}
				}
			}()*/

		} else {
			processFltCacheSearch(subRequest, resultsMap)
		}

	}

	log.Println(" Length of Results Maps : ", len(resultsMap))

	return models.FlightCacheServiceResponse{
		Query:   request,
		Results: resultsMap,
	}

}

func translateRequest(query *models.TfmSearchQuery) *models.SearchRequest {
	return &models.SearchRequest{
		Cached:               false,
		AirlineCode:          "KL", //TODO limited to AF for the poc
		DepartureAirportCode: query.Origin,
		ArrivalAirportCode:   query.Destination,
		DepartureDateTime:    query.DepDate,
		ArrivalDateTime:      "",
		RoundTrip:            isRoundTripJourney(query.JourneyType),
		BookingTime:          time.Now(), //TODO need to check and get it from the request
	}
}

func convertDate(date string) time.Time {
	layout := "2006-01-02"
	t, err := time.Parse(layout, date)

	if err != nil {
		fmt.Println(err)
	}

	return t
}

func isRoundTripJourney(journey string) bool {
	isRoundTripJourney := false
	if journey != "ONEWAY" {
		isRoundTripJourney = true
	}

	return isRoundTripJourney
}

func splitRequestIntoUnits(ws *models.WideSearchQuery) []*models.SearchRequest {

	var ff []*models.SearchRequest
	var dd *models.SearchRequest
	JourneyType := ws.JourneyType
	for _, departureDate := range ws.DepartureDates {
		for _, originAirportCode := range ws.OriginAirportCodes {
			for _, airlineCode := range ws.AirlineCodes {
				for _, arrivalDate := range ws.ArrivalDates {
					for _, source := range ws.Sources {
						for _, destinationAirportCode := range ws.DestinationAirportCodes {

							dd = &models.SearchRequest{
								Cached:               false,
								AirlineCode:          airlineCode,
								DepartureAirportCode: originAirportCode,
								ArrivalAirportCode:   destinationAirportCode,
								DepartureDateTime:    departureDate,
								ArrivalDateTime:      arrivalDate,
								RoundTrip:            isRoundTripJourney(JourneyType),
								BookingTime:          time.Now(),
								Source:               source,
							}
							//fmt.Println(r1," ", r2, " " ,r3, "", r4, " ", r5, " ", r6)
							ff = append(ff, dd)
							//return []models.SearchRequest

						}
					}
				}
			}
		}
	}
	return ff
}

func processSearchQueryRequest(searchRequest *models.SearchRequest, flightCacheProperties *properties.Properties) *models.SearchResult {
	var searchResult *models.SearchResult
	//TODO remove temp initializers after configuration is added
	rec, _ := flightCacheProperties.Get("isFlightCacheRuleConfigurationCheckEnabled")
	isFlightCacheRuleConfigurationCheckEnabled, _ := strconv.ParseBool(rec)
	//If rule configuration check is enabled
	if isFlightCacheRuleConfigurationCheckEnabled {
		//Call to flight cache rule Engine service

		log.Println("Calling Rule Engine ")

		response := ruleEngine.RuleEngineClientResponse(searchRequest)

		if response.Cacheable {
			//call to perform query search
			//if entry present in the cache -> build the response with the value from cache
			// and set required to cache as false
			//else set empty result, required to cache result as true
			searchResult, _, _ = performSearch(searchRequest, flightCacheProperties)
		} else {
			//set empty result, required to cache result as false
			searchResult = &models.SearchResult{
				Result: models.Result{},
				AdditionalInfo: models.AdditionalInfo{
					NeedToBeCached: false,
					ResultType:     "",
				},
			}
		}
	} else {
		//call to perform query search
		//if entry present in the cache -> build the response with the value from cache
		// and set required to cache as false
		//else set empty result, required to cache result as true
		searchResult, _, _ = performSearch(searchRequest, flightCacheProperties)
	}

	return searchResult
}

func performSearchThroughChan(request *models.SearchRequest, p *properties.Properties, searchResult chan *models.SearchResult, searchKey chan *string) {
	//var searchResult *models.SearchResult

	cacheEntryKey := models.DeriveCacheKeyFromRequest(request)
	cacheEntry := redis.Query(cacheEntryKey, p)
	log.Println("Cache entry retrieved: ", cacheEntry.Value)
	searchKey <- &cacheEntryKey
	if cacheEntry.Value == "" {
		//call search service
		searchResult <- &models.SearchResult{
			Result: models.Result{},
			AdditionalInfo: models.AdditionalInfo{
				NeedToBeCached: true,
			},
		}

	} else {
		result, err := service.LoadResultFromCache(cacheEntry.Value)
		if err != nil {
			log.Panicln("Couldn't load result ", err.Error())
		} else {
			searchResult <- &models.SearchResult{
				Result: *result,
				AdditionalInfo: models.AdditionalInfo{
					NeedToBeCached: true,
				},
			}
		}
	}
}

func performSearch(request *models.SearchRequest, p *properties.Properties) (*models.SearchResult, *string, error) {
	var searchResult *models.SearchResult

	cacheEntryKey := models.DeriveCacheKeyFromRequest(request)
	cacheEntry := redis.Query(cacheEntryKey, p)
	fmt.Println("Cache entry retrieved:", cacheEntry.Value, "#")
	if cacheEntry.Value == "" {
		//call search service
		searchResult = &models.SearchResult{
			Result: models.Result{},
			AdditionalInfo: models.AdditionalInfo{
				NeedToBeCached: true,
			},
		}

	} else {
		result, err := service.LoadResultFromCache(cacheEntry.Value)
		if err != nil {
			log.Panicln("Couldn't load result ", err.Error())
		} else {
			searchResult = &models.SearchResult{
				Result: *result,
				AdditionalInfo: models.AdditionalInfo{
					NeedToBeCached: false,
				},
			}
		}
	}

	return searchResult, &cacheEntryKey, nil
}
