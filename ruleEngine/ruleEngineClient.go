package ruleEngine

import (
	"awesomeProject1/models"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	lambdaService "github.com/aws/aws-sdk-go/service/lambda"
	"log"
	"strconv"
	"time"
)

func RuleEngineClientResponse(searchRequest *models.SearchRequest) *models.SearchResponseFromRuleEngine {

	reqBody := models.FlightCacheSearchQuery{
		DepartureDateTimeInUtc: searchRequest.DepartureDateTime,
		ArrivalDateTimeInUTC:   searchRequest.ArrivalDateTime,
		AirlineCode:            searchRequest.AirlineCode,
		BookingTimeInUtc:       dateFmt(searchRequest.BookingTime),
		Origin:                 searchRequest.DepartureAirportCode,
		Destination:            searchRequest.ArrivalAirportCode,
		JourneyType:            journeyType(searchRequest.RoundTrip),
	}
	jsonConversion, err := json.Marshal(reqBody)

	if err != nil {
		log.Fatalln(err)
	}

	//requestBody := strings.NewReader(string(jsonConversion))
	requestBody := string(jsonConversion)

	log.Println("Calling Lambda Service")

	data, err := invokeRuleEngineLambda(requestBody)

	s, err := strconv.Unquote(string(data))

	if err != nil {
		fmt.Println(err, "Error During String unquote")
	}

	//reqURL := "https://gaihazc5ue.execute-api.us-east-2.amazonaws.com/dev/flightCacheConfigServiceLambda"
	//reqURL := "https://gaihazc5ue.execute-api.us-east-2.amazonaws.com/dev/flightCacheLambda"
	//req, err := http.NewRequest(http.MethodPost, reqURL, requestBody)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//res, err := http.DefaultClient.Do(req)
	//
	//log.Println("")
	//
	//if err != nil {
	//	log.Fatalln(err)
	//}

	//data, err := ioutil.ReadAll(res.Body)
	//
	//if err != nil {
	//	panic(err)
	//}
	//res.Body.Close()

	fmt.Println("string converted using unquote : ", s)
	fmt.Println("Original data: ", data)

	resBody := &models.SearchResponseFromRuleEngine{}
	err = json.Unmarshal([]byte(s), &resBody)

	if err != nil {
		panic(err)
	}

	fmt.Println("Final Result from Rule Engine : ", &resBody.Cacheable)

	return resBody
}

func RuleEngineClientResponseThroughChan(searchRequest *models.SearchRequest, response chan *models.SearchResponseFromRuleEngine) {

	reqBody := models.FlightCacheSearchQuery{
		DepartureDateTimeInUtc: searchRequest.DepartureDateTime,
		AirlineCode:            searchRequest.AirlineCode,
		BookingTimeInUtc:       dateFmt(searchRequest.BookingTime),
		Origin:                 searchRequest.DepartureAirportCode,
		Destination:            searchRequest.ArrivalAirportCode,
		JourneyType:            journeyType(searchRequest.RoundTrip),
	}
	jsonConversion, err := json.Marshal(reqBody)

	if err != nil {
		log.Fatalln(err)
	}

	requestBody := string(jsonConversion)

	log.Println("Calling Lambda Service")

	data, err := invokeRuleEngineLambda(requestBody)

	s, err := strconv.Unquote(string(data))

	if err != nil {
		fmt.Println(err, "Error During String unquote")
	}

	log.Println("string converted using unquote : ", s)
	log.Println("Original data: ", data)

	resBody := &models.SearchResponseFromRuleEngine{}
	err = json.Unmarshal([]byte(s), &resBody)

	if err != nil {
		panic(err)
	}

	response <- resBody
}

func journeyType(journey bool) string {

	if !journey {
		return "ONEWAY"
	} else {
		return "RoundTrip"
	}

}

func dateFmt(dateToBeConv time.Time) string {
	y, m, d := dateToBeConv.Date()
	depDate := fmt.Sprintf("%v-%v-%v", y, int(m), d)
	return depDate
}

//func invokeRuleEngineLambda(request string) (response []byte, err error) {
//	sess := session.Must(session.NewSessionWithOptions(session.Options{
//		SharedConfigState: session.SharedConfigEnable,
//	}))
//	log.Println("Session created", sess.Config.Region)
//	client := lambdaService.New(sess, &aws.Config{Region: aws.String("us-east-2")})
//	log.Println("Lambda client created", client.ServiceName)
//	payload := []byte(request)
//	log.Println("Invoking lambda with payload", request)
//	result, err := client.Invoke(&lambdaService.InvokeInput{FunctionName: aws.String("flightCacheConfigServiceLambda"), Payload: payload})
//	if err != nil {
//		fmt.Println("Error calling flightCacheConfigServiceLambda Function ")
//		os.Exit(0)
//	}
//
//	return result.Payload, nil
//}

func invokeRuleEngineLambda(request string) (response []byte, err error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	var (
		invocationType    = "RequestResponse"
		lambdaFunctionARN = "arn:aws:lambda:us-east-2:213716389255:function:flightCacheConfigServiceLambda"
	)
	client := lambdaService.New(sess, &aws.Config{Region: aws.String("us-east-2")})

	fmt.Println("Payload Request: ", request)

	payload := []byte(request)
	log.Println("Calling lambda:", lambdaFunctionARN, " with payload ", request)
	result, err := client.Invoke(&lambdaService.InvokeInput{
		FunctionName:   &lambdaFunctionARN,
		Payload:        payload,
		InvocationType: &invocationType})
	if err != nil {
		log.Println("Error calling flightCacheConfigServiceLambda ", err.Error())
	}

	log.Println("Result.String :", result.String())
	log.Println("Result : ", result)
	log.Println("Payload: String Con : ", string(result.Payload))

	return result.Payload, err
}
