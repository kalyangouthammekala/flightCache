package test

import (
	"fmt"
	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	rule2 = `
rule AgeNameCheck "test" {
    when
      RuleInfo.GetStringLength("9999") > 0  &&
      RuleInfo.Result == ""
    then
      RuleInfo.Result = "String len above 0";
}
`

	rule3 = `
rule AgeNameCheck "test"  salience 10{
    when
      RuleInfo.Compare(User.Name, "Calo")  
    then
      User.Name = "Success";
      Log(User.Name);
      Retract("AgeNameCheck");
}
`
)

// MyRuleInfo serve as example plain Plai Old Go Object.
type MyRuleInfo struct {
	Result string
}

// GetStringLength will return the length of provided string argument
func (p *MyRuleInfo) GetStringLength(sarg string) int {
	return len(sarg)
}

// Compare will compare the equality between the two string.
func (p *MyRuleInfo) Compare(t1, t2 string) bool {
	fmt.Println(t1, t2)
	return t1 == t2
}

// User is an example user struct.
type User struct {
	Name string
	Age  int
	Male bool
}

func TestMyRuleInfo_GetStringLength(t *testing.T) {
	dataContext := ast.NewDataContext()
	RuleInfo := &MyRuleInfo{}
	err := dataContext.Add("RuleInfo", RuleInfo)
	if err != nil {
		t.Fatal(err)
	}

	lib := ast.NewKnowledgeLibrary()
	ruleBuilder := builder.NewRuleBuilder(lib)
	err = ruleBuilder.BuildRuleFromResource("Test", "0.1.1", pkg.NewBytesResource([]byte(rule2)))
	assert.NoError(t, err)
	kb := lib.NewKnowledgeBaseInstance("Test", "0.1.1")
	eng1 := &engine.GruleEngine{MaxCycle: 1}
	err = eng1.Execute(dataContext, kb)
	assert.NoError(t, err)
	assert.Equal(t, "String len above 0", RuleInfo.Result)
	fmt.Println(RuleInfo.Result)
}

func TestMyRuleInfo_Compare(t *testing.T) {
	user := &User{
		Name: "Calo",
		Age:  0,
		Male: true,
	}

	dataContext := ast.NewDataContext()
	err := dataContext.Add("User", user)
	if err != nil {
		t.Fatal(err)
	}
	err = dataContext.Add("RuleInfo", &MyRuleInfo{})
	if err != nil {
		t.Fatal(err)
	}

	lib := ast.NewKnowledgeLibrary()
	ruleBuilder := builder.NewRuleBuilder(lib)

	err = ruleBuilder.BuildRuleFromResource("Test", "0.1.1", pkg.NewBytesResource([]byte(rule3)))
	assert.NoError(t, err)
	kb := lib.NewKnowledgeBaseInstance("Test", "0.1.1")
	eng1 := &engine.GruleEngine{MaxCycle: 100}
	err = eng1.Execute(dataContext, kb)
	assert.NoError(t, err)
	t.Log(user)
	assert.Equal(t, "Success", user.Name, "User should have changed name by rule to Success, but %s", user.Name)
	fmt.Println("user name " + user.Name)
}
