package models

import (
	"fmt"
	"log"
	"strings"
	"time"
)

func (p *SearchResponse) AddDays(inputTime time.Time, days int64) time.Time {
	fmt.Println("adding days ", days)
	return inputTime.AddDate(0, 0, int(days))
}

//func TransformDate(date time.Time) string {
//	dateForKey := strconv.Itoa(date.Year()) + strconv.Itoa(int(date.Month())) + strconv.Itoa(date.Day())
//	fmt.Println("date ", dateForKey)
//	return dateForKey
//}

func DeriveCacheKeyFromRequest(request *SearchRequest) string {

	var (
		key, journeyType string
		keyBuilder       strings.Builder
	)

	keyBuilder.WriteString(request.DepartureAirportCode + KEY_DELIMITER + request.ArrivalAirportCode + KEY_DELIMITER +
		request.AirlineCode + KEY_DELIMITER + request.Source + KEY_DELIMITER +
		request.DepartureDateTime + KEY_DELIMITER)

	//key := ""
	//journeyType := ""
	//arrivalTime := ""
	//JFKHAJ-R-20122021-26122021
	if request.RoundTrip {
		journeyType = "R"
		keyBuilder.WriteString(journeyType + KEY_DELIMITER + request.ArrivalDateTime)
	} else {
		journeyType = "O"
		keyBuilder.WriteString(journeyType)
	}
	key = keyBuilder.String()
	log.Println("Key for Cache ", key)

	return key
}
