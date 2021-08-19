package main

import (
	"awesomeProject1/server"
	"fmt"
)

func main() {
	fmt.Println("Flight Cache app !!")
	/*searchRequest := &models.SearchRequest{
		AirlineCode:          "KL",
		DepartureAirportCode: "AMS",
		ArrivalAirportCode:   "NYC",
		DepartureDateTime: time.Date(
			2021,
			8,
			22,
			12,
			35,
			0,
			0,
			time.UTC,
		),
		ArrivalDateTime: time.Date(
			2021,
			8,
			24,
			12,
			35,
			0,
			0,
			time.UTC,
		),
		RoundTrip:            true,
		BookingTime:          time.Now(),
	}

	ruleEngineResponse := ruleEngine.Execute(searchRequest)
	fmt.Println("Result ", ruleEngineResponse.FromCache)*/

	server.StartServer()
}
