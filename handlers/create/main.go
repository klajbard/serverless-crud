package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
	"github.com/klajbard/serverless-crud/global"
)

// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

func create(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if strings.TrimSpace(event.Body) == "" {
		return global.ErrResponse(http.StatusBadRequest, "empty request body"), nil
	}

	animal := global.Animal{}

	if err := json.Unmarshal([]byte(event.Body), &animal); err != nil {
		fmt.Println(err)
		return global.ErrResponse(http.StatusBadRequest, "failed to parse request body"), nil
	}

	id := uuid.New()
	animal.Id = id.String()

	item, err := dynamodbattribute.MarshalMap(&animal)
	if err != nil {
		fmt.Println(err)
		return global.ErrResponse(http.StatusInternalServerError, "unable to marshal animal data"), nil
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(global.TableName),
		Item:      item,
	}

	_, err = global.DB.PutItemWithContext(ctx, input)

	if err != nil {
		fmt.Println(err)
		return global.ErrResponse(http.StatusInternalServerError, "failed to create animal"), nil
	}

	return global.Response(http.StatusOK, animal.Id), nil
}

func main() {
	lambda.Start(create)
}
