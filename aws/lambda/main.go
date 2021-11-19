package main

import (
	"awesomeProject1/config"
	"awesomeProject1/controller"
	"encoding/json"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/magiconair/properties"
	"strings"

	/*"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"*/
	"log"
	"net/http"
)

var (
	flightCacheProperties *properties.Properties
)

func init() {
	flightCacheProperties = config.LoadProperties()
}

func main() {
	lambda.Start(HandleRequest)
}

func HandleRequest(input interface{}) (interface{}, error) {

	var err error
	inputRequest, err := json.Marshal(input)
	if err != nil {
		log.Println("unable to marshal input to json")
	}

	inputRequestAsString := string(inputRequest)
	containsHeader := strings.Contains(inputRequestAsString, "headers")
	log.Println("does the request contain header ", containsHeader)
	if !containsHeader {
		log.Println("calling lambda for request ", inputRequestAsString)
		//apiResponse, err := performRuleCheck(inputRequestAsString)
		response, err := controller.SearchHandler(inputRequestAsString, flightCacheProperties)
		if err != nil {
			return err.Error(), err
		} else {
			return string(response), err
		}
	} else if containsHeader {
		var apiRequest events.APIGatewayProxyRequest

		err = json.Unmarshal(inputRequest, &apiRequest)
		if err != nil {
			log.Panic(err.Error())
		} else {
			log.Println("Entering into API Request Handler: ")
			return handleAPIRequest(apiRequest)
		}

	} else {
		err = errors.New("nothing executed")
	}

	if containsHeader {
		return events.APIGatewayProxyResponse{Body: "Nothing Executed", StatusCode: http.StatusOK, IsBase64Encoded: false, Headers: nil}, err
	} else {
		return "Nothing Executed", err
	}

}

func handleAPIRequest(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	log.Println("Path:", req.Path)
	if req.HTTPMethod == http.MethodPost && req.Path == "/flightCache/search" {
		log.Println("Inside POST")
		log.Println("Body ", req.Body)

		res, err := controller.SearchHandler(req.Body, flightCacheProperties)

		return events.APIGatewayProxyResponse{Body: string(res), StatusCode: 200}, err
	} else {
		log.Println("Method : ", req.Headers["httpMethod"])
		log.Println("Body ", req.Body)
		err := errors.New("invalid request")
		return events.APIGatewayProxyResponse{Body: "Method not ok", StatusCode: 502}, err
	}

}
