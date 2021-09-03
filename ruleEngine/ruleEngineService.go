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
	"strings"
)

func LoadRules(name, version string) *ast.KnowledgeBase {
	lib := ast.NewKnowledgeLibrary()
	ruleBuilder := builder.NewRuleBuilder(lib)
	rulesFile, err := filepath.Abs("../resources/rules.json")
	//rulesFile, err := filepath.Abs("../resources/rules.json") //TODO change when tested from controller
	//TODO need to solve by having a uniform access from test file as well as controller
	fmt.Println(rulesFile)

	jsonData, err := ioutil.ReadFile(rulesFile)
	if err != nil {
		if strings.Contains(err.Error(), "The system cannot find the path specified") {
			err = nil
			rulesFile, err = filepath.Abs("../awesomeProject1/resources/rules.json")
			jsonData, err = ioutil.ReadFile(rulesFile)
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
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
	err = dataContext.Add("RuleInfo", searchResponse)
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
	//fetch matching rules
	ruleEntries, err := eng1.FetchMatchingRules(*dataContext, kb)
	if err != nil {
		fmt.Println("Unable to fetch all matching rules")
		panic(err)
	} else {
		fmt.Println("Matching Rule(s):")
		for i := 0; i < len(ruleEntries); i++ {
			fmt.Println(ruleEntries[i].RuleName + " : " + ruleEntries[i].RuleDescription)
		}
	}

	err = eng1.Execute(*dataContext, kb)
	if err != nil {
		fmt.Println(err)
	}
	return searchResponse
}
