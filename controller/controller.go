package controller

import (
	"awesomeProject1/models"
	"awesomeProject1/searchService"
	"awesomeProject1/service"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func StartServer() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/search", searchHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, %s!", r.URL.Path[1:])
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	tfmQuery := models.TfmSearchQuery{}
	err = json.Unmarshal([]byte(body), &tfmQuery)
	if err != nil {
		panic(err)
	}

	searchRequest := translateRequest(&tfmQuery)
	flightCacheService := &service.FlightCacheService{
		Request: searchRequest,
		Response: &models.SearchResponse{
			FromCache:            false,
			AirlineCode:          searchRequest.AirlineCode,
			DepartureAirportCode: searchRequest.DepartureAirportCode,
			ArrivalAirportCode:   searchRequest.ArrivalAirportCode,
			RoundTrip:            searchRequest.RoundTrip,
			BookingTime:          searchRequest.BookingTime,
		},
		SearchService: &searchService.DummySearchServiceImpl{}, //TODO replace with actual implementation before deployment
	}
	kbDetails := &models.KnowledgeBaseForCacheRule{
		Name:    "Test",
		Version: "0.0.1",
	}
	response := flightCacheService.Search(kbDetails)
	//searchResponse := ruleEngine.Execute(searchRequest, nil, "", "")
	fmt.Println(tfmQuery.Destination, " ", response.FromCache)
	responseData, err := json.Marshal(response.TfmRessponse)
	if err != nil {
		panic(err)
	}
	w.Write(responseData)
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
