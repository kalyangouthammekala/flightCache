package controller

import (
	"awesomeProject1/models"
	"awesomeProject1/searchService"
	"awesomeProject1/service"
	"encoding/json"
	"fmt"
	"github.com/magiconair/properties"
	"log"
	"net/http"
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
	response := flightCacheService.Search(kbDetails, flightCacheProperties)
	//searchResponse := ruleEngine.Execute(searchRequest, nil, "", "")
	fmt.Println(tfmQuery.Destination, " ", response.FromCache)
	responseData, err := json.Marshal(response) //(response.TfmRessponse)
	if err != nil {
		panic(err)
	}
	//_, err = w.Write(responseData)
	//if err != nil {
	//	fmt.Println("Error in writing response: ", err.Error())
	//}
	return responseData, err
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
