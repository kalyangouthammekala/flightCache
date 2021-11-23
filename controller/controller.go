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
func SearchHandler(r string, flightCacheProperties *properties.Properties) ([]byte, error) {
	log.Println("Entering into SearchHandler")
	//body, err := ioutil.ReadAll(r.Body)
	//if err != nil {
	//	panic(err)
	//}
	tfmQuery := models.TfmSearchQuery{}
	err := json.Unmarshal([]byte(r), &tfmQuery)
	if err != nil {
		panic(err)
	}

	// wideSearch as json interface
	// loop through the combinations
	// and return search objects.

	// FOR WIDE SEARCH

	ws := models.WideSearchQuery{}
	err = json.Unmarshal([]byte(r), &ws)
	if err != nil {
		panic(err)
	}

	searchRequestArray := searchRequestArrayObj(&ws)

	for i, request := range searchRequestArray {
		// call the search service for the i times
		fmt.Println(i, "Search Objects")
		_ = processSearchQueryRequest(request, flightCacheProperties)
		if err != nil {
			log.Panic(err)
		}
	}

	//searchRequest := translateRequest(&tfmQuery)

	//_, err = w.Write(responseData)
	//if err != nil {
	//	fmt.Println("Error in writing response: ", err.Error())
	//}
	return []byte("test"), err
}

func translateRequest(query *models.TfmSearchQuery) *models.SearchRequest {
	return &models.SearchRequest{
		Cached:               false,
		AirlineCode:          "KL", //TODO limited to AF for the poc
		DepartureAirportCode: query.Origin,
		ArrivalAirportCode:   query.Destination,
		DepartureDateTime:    convertDate(query.DepDate),
		ArrivalDateTime:      time.Time{},
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

func searchRequestArrayObj(ws *models.WideSearchQuery) []*models.SearchRequest {

	var ff []*models.SearchRequest
	var dd *models.SearchRequest
	JourneyType := ws.JourneyType
	for _, r1 := range ws.DepartureDates {
		for _, r2 := range ws.OriginAirportCodes {
			for _, r3 := range ws.AirlineCodes {
				for _, r4 := range ws.ArrivalDates {
					for _, _ = range ws.Sources {
						for _, r6 := range ws.DestinationAirportCodes {

							dd = &models.SearchRequest{
								Cached:               false,
								AirlineCode:          r3,
								DepartureAirportCode: r2,
								ArrivalAirportCode:   r6,
								DepartureDateTime:    convertDate(r1),
								ArrivalDateTime:      convertDate(r4),
								RoundTrip:            isRoundTripJourney(JourneyType),
								BookingTime:          time.Now(),
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

		response := ruleEngine.RuleEngineClientResponse(searchRequest)
		if response.Cacheable {
			//call to perform query search
			//if entry present in the cache -> build the response with the value from cache
			// and set required to cache as false
			//else set empty result, required to cache result as true
			searchResult = performSearch(searchRequest, flightCacheProperties)
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
		searchResult = performSearch(searchRequest, flightCacheProperties)
	}

	return searchResult
}

func performSearch(request *models.SearchRequest, p *properties.Properties) *models.SearchResult {
	var searchResult *models.SearchResult

	cacheEntryKey := models.DeriveCacheKeyFromRequest(request)
	cacheEntry := redis.Query(cacheEntryKey, p)
	fmt.Println("Cache entry retrieved: ", cacheEntry.Value)
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
					NeedToBeCached: true,
				},
			}
		}
	}

	return searchResult
}
