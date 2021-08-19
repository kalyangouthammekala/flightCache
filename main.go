package main

import (
	"awesomeProject1/models"
	"awesomeProject1/ruleEngine"
	"fmt"
	"time"
)

func main() {
	fmt.Println("Back to golang")
	searchRequest := &models.SearchRequest{
		AirlineCode:          "KL",
		DepartureAirportCode: "AMS",
		ArrivalAirportCode:   "NYC",
		RoundTrip:            true,
		BookingTime:          time.Now(),
	}

	ruleEngineResponse := ruleEngine.Execute(searchRequest)
	fmt.Println("Result ", ruleEngineResponse.FromCache)
}
