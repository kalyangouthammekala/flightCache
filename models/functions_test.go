package models

import (
	"fmt"
	"testing"
	"time"
)

var oneWayResult = SearchRequest{
	Cached:               false,
	AirlineCode:          "KL",
	DepartureAirportCode: "AMS",
	ArrivalAirportCode:   "HAJ",
	DepartureDateTime:    "2021-11-22",
	ArrivalDateTime:      "",
	RoundTrip:            false,
	BookingTime:          time.Time{},
	Source:               "NDC",
}

var roundTripResult = SearchRequest{
	Cached:               false,
	AirlineCode:          "KL",
	DepartureAirportCode: "AMS",
	ArrivalAirportCode:   "HAJ",
	DepartureDateTime:    "2021-11-22",
	ArrivalDateTime:      "2021-11-23",
	RoundTrip:            true,
	BookingTime:          time.Time{},
	Source:               "NDC",
}

func TestDeriveCacheKeyFromRequestONEWAY(t *testing.T) {
	var (
		keyO string
	)

	keyO = DeriveCacheKeyFromRequest(&oneWayResult)

	fmt.Println(keyO)
}

func TestDeriveCacheKeyFromRequestRound(t *testing.T) {
	var (
		key string
	)

	key = DeriveCacheKeyFromRequest(&roundTripResult)
	fmt.Println(key)
}
