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
	TfmRessponse         TfmResponse
	BookingTime          time.Time
}

type TfmResponse struct {
	Query            map[string]interface{}   `json:"query"`
	Routes           map[string]Route         `json:"routes"`
	Segments         map[string]Segment       `json:"segments"`
	Combinations     []Combination            `json:"combinations"`
	Ancillaries      []Ancillary              `json:"ancillaries"`
	AdditionalParams map[string]string        `json:"additionalParams,omitempty"`
	ResponseTimes    map[string]time.Duration `json:"responseTimes,omitempty"`
}
type Route struct {
	Id                       string            `json:"id"`
	Stops                    int8              `json:"stops"`
	ElapsedFlyingTimeMinutes int               `json:"elapsedFlyingTimeMinutes"`
	SegmentIDs               []string          `json:"segmentIDs"`
	AdditionalParams         map[string]string `json:"additionalParams,omitempty"`
}

type Segment struct {
	Id                  string            `json:"id"`
	Origin              string            `json:"origin"`
	OriginTerminal      string            `json:"originTerminal"`
	Destination         string            `json:"destination"`
	DestinationTerminal string            `json:"destinationTerminal"`
	DepartureTime       string            `json:"departureTime"`
	ArrivalTime         string            `json:"arrivalTime"`
	FlightNumber        string            `json:"flightNumber"`
	AirplaneType        string            `json:"airplaneType"`
	MarketingCarrier    string            `json:"marketingCarrier"`
	OperationCarrier    string            `json:"operatingCarrier"`
	AdditionalParams    map[string]string `json:"additionalParams,omitempty"`
}

type Combination struct {
	TotalFareAmount  float64           `json:"totalFareAmount"`
	TotalTaxAmount   float64           `json:"totalTaxAmount"`
	Fares            []TfmFare         `json:"fares"`
	RouteIDs         []string          `json:"routeIDs"`
	TariffType       string            `json:"tariffType"`
	AdditionalParams map[string]string `json:"additionalParams,omitempty"`
}

type TfmFare struct {
	PaxId        string        `json:"paxId"`
	PaxType      string        `json:"paxType"`
	FareAmount   float64       `json:"fareAmount"`
	TaxAmount    float64       `json:"taxAmount"`
	FareProducts []FareProduct `json:"fareProducts"`
	Vcc          string        `json:"vcc"`
}

type FareProduct struct {
	SegmentID        string            `json:"segmentID"`
	CabinProduct     string            `json:"cabinProduct"`
	FareBase         string            `json:"fareBase"`
	AncillaryIDs     []string          `json:"ancillaryIDs"`
	AdditionalParams map[string]string `json:"additionalParams,omitempty"`
}

type Ancillary struct {
	Id               string            `json:"id"`
	Type             string            `json:"type"`
	AdditionalParams map[string]string `json:"additionalParams,omitempty"`
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
