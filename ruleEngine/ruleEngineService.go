package ruleEngine

import (
	"awesomeProject1/models"
	"fmt"
	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
	"io/ioutil"
	"path/filepath"
)

func LoadRules(name, version string) *ast.KnowledgeBase {
	lib := ast.NewKnowledgeLibrary()
	ruleBuilder := builder.NewRuleBuilder(lib)
	rulesFile, err := filepath.Abs("../awesomeProject1/resources/rules.json")
	fmt.Println(rulesFile)

	jsonData, err := ioutil.ReadFile(rulesFile)
	if err != nil {
		panic(err)
	}
	ruleset, err := pkg.ParseJSONRuleset(jsonData)
	if err != nil {
		panic(err)
	}

	fmt.Println("Parsed ruleset: ")
	fmt.Println(ruleset)
	err = ruleBuilder.BuildRuleFromResource(name, version, pkg.NewBytesResource([]byte(ruleset)))
	kb := lib.NewKnowledgeBaseInstance(name, version)
	return kb
}

func LoadDataContextToKnowledgeBase(searchRequest *models.SearchRequest,
	searchResponse *models.SearchResponse) *ast.IDataContext {
	dataContext := ast.NewDataContext()
	err := dataContext.Add("FltSearchRequest", searchRequest)
	if err != nil {
		fmt.Println("Error while loading search request (fact): ", err)
	}
	*searchResponse = models.SearchResponse{
		FromCache:            false,
		AirlineCode:          searchRequest.AirlineCode,
		DepartureAirportCode: searchRequest.DepartureAirportCode,
		ArrivalAirportCode:   searchRequest.ArrivalAirportCode,
		RoundTrip:            searchRequest.RoundTrip,
		BookingTime:          searchRequest.BookingTime,
	}
	err = dataContext.Add("Pogo", searchResponse)
	if err != nil {
		fmt.Println(err)
	}
	return &dataContext
}

func Execute(searchRequest *models.SearchRequest, searchResponse *models.SearchResponse,
	name, version string) *models.SearchResponse {
	dataContext := LoadDataContextToKnowledgeBase(searchRequest, searchResponse)
	kb := LoadRules(name, version)
	eng1 := &engine.GruleEngine{MaxCycle: 100}
	err := eng1.Execute(*dataContext, kb)
	if err != nil {
		fmt.Println(err)
	}
	return searchResponse
}

/*func Execute(searchRequest *models.SearchRequest) *models.SearchResponse {
	dataContext := ast.NewDataContext()
	err := dataContext.Add("FltSearchRequest", searchRequest)
	if err != nil {
		fmt.Println("Error while loading search request (fact): ", err)
	}
	result := &models.SearchResponse{
		FromCache:            false,
		AirlineCode:          searchRequest.AirlineCode,
		DepartureAirportCode: searchRequest.DepartureAirportCode,
		ArrivalAirportCode:   searchRequest.ArrivalAirportCode,
		RoundTrip:            searchRequest.RoundTrip,
		BookingTime:          searchRequest.BookingTime,
	}
	err = dataContext.Add("Pogo", result)
	if err != nil {
		fmt.Println(err)
	}

	lib := ast.NewKnowledgeLibrary()
	ruleBuilder := builder.NewRuleBuilder(lib)
	rulesFile, err := filepath.Abs("../awesomeProject1/resources/rules.json")
	fmt.Println(rulesFile)

	jsonData, err := ioutil.ReadFile(rulesFile)
	if err != nil {
		panic(err)
	}
	ruleset, err := pkg.ParseJSONRuleset(jsonData)
	if err != nil {
		panic(err)
	}

	fmt.Println("Parsed ruleset: ")
	fmt.Println(ruleset)
	err = ruleBuilder.BuildRuleFromResource("Test", "0.1.1", pkg.NewBytesResource([]byte(ruleset)))
	kb := lib.NewKnowledgeBaseInstance("Test", "0.1.1")
	eng1 := &engine.GruleEngine{MaxCycle: 100}
	err = eng1.Execute(dataContext, kb)
	if err != nil {
		fmt.Println(err)
	}
	return result
}*/
