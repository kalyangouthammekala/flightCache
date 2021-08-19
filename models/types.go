package models

import "time"

type SearchRequest struct {
	Cached               bool
	AirlineCode          string
	DepartureAirportCode string
	ArrivalAirportCode   string
	RoundTrip            bool
	BookingTime          time.Time
}

type SearchResponse struct {
	FromCache            bool
	AirlineCode          string
	DepartureAirportCode string
	ArrivalAirportCode   string
	RoundTrip            bool
	BookingTime          time.Time
}
