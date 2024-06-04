package handlers

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

func apiResponse(statusCode int, body interface{})(*events.APIGatewayProxyResponse, error){
	resp := events.APIGatewayProxyResponse{Headers: map[string]string{"Content-Type": "application/json"}}
	resp.StatusCode = statusCode

	stringBody, _ := json.Marshal(body)
	resp.Body = string(stringBody)

	return &resp, nil
}