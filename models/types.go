package models

import "time"

type SearchRequest struct {
	Cached               bool
	AirlineCode          string
	DepartureAirportCode string
	ArrivalAirportCode   string
	DepartureDateTime    time.Time
	ArrivalDateTime      time.Time
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

type TfmSearchQuery struct {
	Adult           int    `json:"adult"`
	APICall         string `json:"apiCall"`
	BaggageType     string `json:"baggageType"`
	Channel         string `json:"channel"`
	Child           int    `json:"child"`
	ContextToken    string `json:"contextToken"`
	CorrelationID   string `json:"correlationId"`
	Currency        string `json:"currency"`
	DepDate         string `json:"depDate"`
	Destination     string `json:"destination"`
	ExistingPools   int    `json:"existingPools"`
	Infant          int    `json:"infant"`
	JourneyType     string `json:"journeyType"`
	Origin          string `json:"origin"`
	PoolRequestNr   int    `json:"poolRequestNr"`
	Product         string `json:"product"`
	RequestID       string `json:"requestId"`
	SessionID       string `json:"sessionId"`
	Source          string `json:"source"`
	TrackingEnabled bool   `json:"trackingEnabled"`
	UsePolling      bool   `json:"usePolling"`
}

type CacheEntry struct {
	Key   string
	Value string
}

type KnowledgeBaseForCacheRule struct {
	Name    string
	Version string
}
