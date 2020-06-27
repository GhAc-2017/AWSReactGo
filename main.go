package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/aws/aws-lambda-go/events"
	
	"github.com/aws/aws-lambda-go/lambda"
)
 
var isbnRegexp = regexp.MustCompile(`[0-9]{3}\-[0-9]{10}`)
var errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)

// {
// 	"name": "Test Event",
// 	"description": "A Test Event",
// 	"status": "Idle",
// 	"schedule": {
// 		"start_time": "string",
// 		"stop_time": "string"
// 	}
// }
type schedule struct {
	StartTime string `json:"start_time"`
	StopTime  string `json:"stop_time"`
}
type event struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Status      string            `json:"status"`
	Schedule    map[string]string `json:"schedule"`
	// Schedule    schedule `json:"schedule"`
}

func show() (*event, error) {
	//static item

	// ev := &event{
	// 	Name:        "Test Event",
	// 	Description: "A Test Event",
	// 	Status:      "Idle",
	// 	Schedule:    schedule{StartTime: "sfdfsdf", StopTime: "dfsdf"}}

	//Dynamic req from dynamo db

	ev, err := getItem("Test Event")
	if err != nil {
		return nil, err
	}

	return ev, nil
}

// Add a helper for handling errors. This logs any error to os.Stderr
// and returns a 500 Internal Server Error response that the AWS API
// Gateway understands.
func serverError(err error) (events.APIGatewayProxyResponse, error) {
	errorLogger.Println(err.Error())

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       http.StatusText(http.StatusInternalServerError),
	}, nil
}

// Similarly add a helper for send responses relating to client errors.
func clientError(status int) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       http.StatusText(status),
	}, nil
}

func showfromDB(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Get the `isbn` query string parameter from the request and
	// validate it.
	name := req.QueryStringParameters["isbn"]
	if !isbnRegexp.MatchString(name) {
		return clientError(http.StatusBadRequest)
	} // Fetch the book record from the database based on the isbn value.
	ev, err := getItem(name)
	if err != nil {
		return serverError(err)
	}
	if ev == nil {
		return clientError(http.StatusNotFound)
	}

	// The APIGatewayProxyResponse.Body field needs to be a string, so
	// we marshal the book record into JSON.
	js, err := json.Marshal(ev)
	if err != nil {
		return serverError(err)
	}

	// Return a response with a 200 OK status and the JSON book record
	// as the body.
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(js),
	}, nil

}

func main() {
	lambda.Start(show)
}
