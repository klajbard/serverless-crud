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

type UploadImageRequest struct {
	Id     string `json:"id"`
	Base64 string `json:"base64"`
}

type Animal struct {
	Id          string `json:"id" dynamodbav:"id"`
	Name        string `json:"name" dynamodbav:"name"`
	Description string `json:"description" dynamodbav:"description"`
	Status      bool   `json:"status" dynamodbav:"status"`
	Avatar      string `json:"avatar" dynamodbav:"avatar"`
	Breed       string `json:"breed" dynamodbav:"breed"`
	Birth       int64  `json:"birth" dynamodbav:"birth"`
}

type AnimalUpdate struct {
	Name        string `json:":name" dynamodbav:":name,omitempty"`
	Description string `json:":description" dynamodbav:":description,omitempty"`
	Status      bool   `json:":status" dynamodbav:":status"`
	Avatar      string `json:":avatar" dynamodbav:":avatar,omitempty"`
	Breed       string `json:":breed" dynamodbav:":breed,omitempty"`
	Birth       int64  `json:":birth" dynamodbav:":birth,omitempty"`
}

type SimpleResponse struct {
	Response string `json:"response"`
}

var TableName string

var DB dynamodb.DynamoDB
var Session *session.Session

func init() {
	TableName = os.Getenv("ANIMALS_TABLE_NAME")
	Session = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	DB = *dynamodb.New(Session)
}

func Response(code int, object interface{}) events.APIGatewayProxyResponse {
	marshalled, err := json.Marshal(object)
	if err != nil {
		return ErrResponse(http.StatusInternalServerError, err.Error())
	}

	return events.APIGatewayProxyResponse{
		StatusCode: code,
		Headers: map[string]string{
			"Content-Type":                     "application/json",
			"Access-Control-Allow-Origin":      "http://localhost:5173",
			"Access-Control-Allow-Credentials": "true",
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
			"Content-Type":                     "application/json",
			"Access-Control-Allow-Origin":      "http://localhost:5173",
			"Access-Control-Allow-Credentials": "true",
		},
		Body: string(messageBytes),
	}
}
