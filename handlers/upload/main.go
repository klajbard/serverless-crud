package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/klajbard/serverless-crud/global"
)

// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

func upload(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5000*time.Millisecond)
	defer cancel()

	if strings.TrimSpace(event.Body) == "" {
		return global.ErrResponse(http.StatusBadRequest, "empty request body"), nil
	}

	animal := global.Animal{}

	if err := json.Unmarshal([]byte(event.Body), &animal); err != nil {
		return global.ErrResponse(http.StatusBadRequest, "failed to parse request body"), nil
	}

	item, err := dynamodbattribute.MarshalMap(&animal)
	if err != nil {
		log.Println(err)
		return global.ErrResponse(http.StatusInternalServerError, "unable to marshal animal data"), nil
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(global.TableName),
		Item:      item,
	}

	_, err = global.DB.PutItemWithContext(ctx, input)

	if err != nil {
		log.Println(err)
		return global.ErrResponse(http.StatusInternalServerError, "failed to create animal"), nil
	}

	return global.Response(http.StatusOK, "successfully created"), nil
}

func main() {
	lambda.Start(upload)
}
