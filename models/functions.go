package models

import (
	"fmt"
	"strconv"
	"time"
)

func (p *SearchResponse) AddDays(inputTime time.Time, days int64) time.Time {
	fmt.Println("adding days ", days)
	return inputTime.AddDate(0, 0, int(days))
}

func TransformDate(date time.Time) string {
	dateForKey := strconv.Itoa(date.Day()) + date.Month().String() + strconv.Itoa(date.Year())
	fmt.Println("date ", dateForKey)
	return dateForKey
}

func DeriveCacheKeyFromRequest(request *SearchRequest) string {
	key := ""
	journeyType := ""
	arrivalTime := ""
	//JFKHAJ-R-20122021-26122021
	if request.RoundTrip {
		journeyType = "R"
		arrivalTime = "-" + TransformDate(request.ArrivalDateTime)
	} else {
		journeyType = "O"
	}
	key = request.DepartureAirportCode + request.ArrivalAirportCode + "-" + journeyType + "" +
		"-" + TransformDate(request.DepartureDateTime) +
		arrivalTime

	fmt.Println("Key for Cache ", key)

	return key
}
