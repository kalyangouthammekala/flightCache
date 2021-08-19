package main

import (
	"fmt"
	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
	_ "github.com/hyperjumptech/grule-rule-engine/pkg"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const ruleKLAMSNYCRoundTrip3To7DaysOfBookingInJson = `[{
  "name": "FltCacheCheck",
  "desc": "when airline is KLM and flt departure is within 3 to 7 days of booking",
  "salience": 10,
  "when": "Pogo.AirlineCode == \"KL\" && FltSearchResult.DepartureTime > Pogo.AddDays(Pogo.BookingTime, 3) && FltSearchResult.DepartureTime < Pogo.AddDays(Pogo.BookingTime, 7) && Pogo.Result == false",
  "then": [
    "Pogo.Result = true",
    "Log(\"Result could be cached\")"
  ]
},
{
  "name": "FltCacheCheck1",
  "desc": "when airline is KLM",
  "salience": 10,
  "when": "Pogo.AirlineCode == \"KL\" && Pogo.Result == false",
  "then": [
    "Pogo.Result = true",
    "Log(\"Result could be cached\")"
  ]
}]`

const arrayJSONData = `[{
    "name": "SpeedUp",
    "desc": "When testcar is speeding up we keep increase the speed.",
    "salience": 10,
    "when": "TestCar.SpeedUp == true && TestCar.Speed < TestCar.MaxSpeed",
    "then": [
        "TestCar.Speed = TestCar.Speed + TestCar.SpeedIncrement",
        "DistanceRecord.TotalDistance = DistanceRecord.TotalDistance + TestCar.Speed",
        "Log(\"Speed increased\")"
    ]
}]`

const (
	ruleKLAMSNYCRoundTrip3To7DaysOfBooking = `
{
  "name": "FltCacheCheck",
  "desc": "when airline is KLM and flt departure is within 3 to 7 days of booking",
  "salience": 10,
  "when": "Pogo.AirlineCode == \"KL\" && FltSearchResult.DepartureTime > Pogo.AddDays(Pogo.BookingTime, 3) && FltSearchResult.DepartureTime < Pogo.AddDays(Pogo.BookingTime, 7) &&,,Pogo.Result == false",
  "then": [
    "Pogo.Result = true",
    "Log(\"Result could be cached\")"
  ]
}
`
)

// CacheResult serve as example plain Plai Old Go Object.
type CacheResult struct {
	Result               bool
	AirlineCode          string
	DepartureAirlineCode string
	ArrivalAirlineCode   string
	RoundTrip            bool
	BookingTime          time.Time
}

// GetStringLength will return the length of provided string argument
/*func (p *CacheResult) ShouldCache(bookingTime time.Time, searchResult FlightSearchResult) int {
	return p.IsWithinCacheDateRange(bookingTime, searchResult.)
}*/

func (p *CacheResult) AddDays(inputTime time.Time, days int64) time.Time {
	fmt.Println("adding days ", days)
	return inputTime.AddDate(0, 0, int(days))
}

func (p *CacheResult) IsWithinCacheDateRange(bookingTime, flightTime time.Time) bool {
	lastDate := bookingTime.AddDate(0, 0, 7)
	startDate := bookingTime.AddDate(0, 0, 3)
	fmt.Println("Is within Cache Date range " + startDate.String() + " and " + lastDate.String())

	return flightTime.Before(lastDate) && flightTime.After(startDate)
}

// User is an example user struct.
type FlightSearchResult struct {
	AirlineCode          string
	DepartureAirlineCode string
	ArrivalAirlineCode   string
	RoundTrip            bool
	DepartureTime        time.Time
	ArrivalTime          time.Time
}

func Test_IsWithinCacheDateRange(t *testing.T) {
	fltSearchResult := &FlightSearchResult{
		AirlineCode:          "KL",
		DepartureAirlineCode: "AMS",
		ArrivalAirlineCode:   "NYC",
		RoundTrip:            true,
		DepartureTime: time.Date(
			2021,
			8,
			22,
			12,
			35,
			0,
			0,
			time.UTC,
		),
		ArrivalTime: time.Date(
			2021,
			8,
			24,
			12,
			35,
			0,
			0,
			time.UTC,
		),
	}

	/*fltSearchResult1 := &FlightSearchResult{
		AirlineCode:          "KL",
		DepartureAirlineCode: "AMS",
		ArrivalAirlineCode:   "NYC",
		RoundTrip:            true,
		DepartureTime:        time.Date(
			2021,
			8,
			18,
			12,
			35,
			0,
			0,
			time.UTC,
		),
		ArrivalTime:         time.Date(
			2021,
			8,
			19,
			12,
			35,
			0,
			0,
			time.UTC,
		),
	}*/

	dataContext := ast.NewDataContext()
	err := dataContext.Add("FltSearchResult", fltSearchResult)
	//err := dataContext.Add("FltSearchResult", fltSearchResult1)
	if err != nil {
		t.Fatal(err)
	}
	result := &CacheResult{
		Result:               false,
		AirlineCode:          "KL",
		DepartureAirlineCode: "AMS",
		ArrivalAirlineCode:   "NYC",
		RoundTrip:            true,
		BookingTime:          time.Now(),
	}
	err = dataContext.Add("Pogo", result)
	if err != nil {
		t.Fatal(err)
	}

	lib := ast.NewKnowledgeLibrary()
	ruleBuilder := builder.NewRuleBuilder(lib)

	//err = ruleBuilder.BuildRuleFromResource("Test", "0.1.1", pkg.NewBytesResource([]byte(ruleAcceptKLMFlightsOnly))) //working
	/*err = ruleBuilder.BuildRuleFromResource("Test", "0.1.1", pkg.NewBytesResource([]byte(ruleKLAMSNYCRoundTrip3To7DaysOfBooking)))
	if err != nil {
		t.Fatal(err)
	}*/
	ruleset, err := pkg.ParseJSONRuleset([]byte(ruleKLAMSNYCRoundTrip3To7DaysOfBookingInJson))
	if err != nil {
		panic(err)
	}
	fmt.Println("Parsed ruleset: ")
	fmt.Println(ruleset)
	err = ruleBuilder.BuildRuleFromResource("Test", "0.1.1", pkg.NewBytesResource([]byte(ruleset)))

	/*//bs := pkg.NewBytesResource([]byte(ruleKLAMSNYCRoundTrip3To7DaysOfBookingInJson))
	err = ruleBuilder.BuildRuleFromResource("Test", "0.1.1", bs)
	if err != nil {
		t.Fatal("Failed to parse rule from json: " + err.Error())
	}

	t.Log(bs)*/
	//assert.NoError(t, err)
	kb := lib.NewKnowledgeBaseInstance("Test", "0.1.1")
	eng1 := &engine.GruleEngine{MaxCycle: 100}
	err = eng1.Execute(dataContext, kb)
	if err != nil {
		t.Fatal(err)
	}
	assert.NoError(t, err)
	t.Log(fltSearchResult.DepartureTime.String())
	assert.Equal(t, true, result.Result, "Cache validation should pass but didn't for flt time", fltSearchResult.DepartureTime, " booking time ", result.BookingTime)
	fmt.Println(" booking time" + result.BookingTime.String())
}
