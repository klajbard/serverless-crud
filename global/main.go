package global

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type AnimalsResults struct {
	Animals []Animal `json:"animals"`
	Next    string   `json:"next"`
}

type Animal struct {
	Id          string   `json:"id" dynamodbav:"id"`
	Name        string   `json:"name" dynamodbav:"name"`
	Description string   `json:"description" dynamodbav:"description"`
	Status      bool     `json:"status" dynamodbav:"status"`
	Images      []string `json:"images" dynamodbav:"images"`
}

type AnimalUpdate struct {
	Id          string   `json:":id" dynamodbav:":id,omitempty"`
	Name        string   `json:":name" dynamodbav:":name,omitempty"`
	Description string   `json:":description" dynamodbav:":description,omitempty"`
	Status      bool     `json:":status" dynamodbav:":status,omitempty"`
	Images      []string `json:":images" dynamodbav:":images,omitempty"`
}

var TableName string

var DB dynamodb.DynamoDB

func init() {
	TableName = os.Getenv("ANIMALS_TABLE_NAME")
	session := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	DB = *dynamodb.New(session)
}

func Response(code int, object interface{}) events.APIGatewayProxyResponse {
	marshalled, err := json.Marshal(object)
	if err != nil {
		return ErrResponse(http.StatusInternalServerError, err.Error())
	}

	return events.APIGatewayProxyResponse{
		StatusCode: code,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body:            string(marshalled),
		IsBase64Encoded: false,
	}
}

func ErrResponse(status int, body string) events.APIGatewayProxyResponse {
	message := map[string]string{
		"message": body,
	}

	messageBytes, _ := json.Marshal(&message)

	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(messageBytes),
	}
}
