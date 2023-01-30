package main

import (
	"context"
	"net/http"
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

func find(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5000*time.Millisecond)
	defer cancel()

	id := event.QueryStringParameters["id"]

	if id == "" {
		return global.ErrResponse(http.StatusBadRequest, "no item specified id to find"), nil
	}

	animal := global.Animal{}

	input := &dynamodb.GetItemInput{
		TableName: aws.String(global.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {S: &id},
		},
	}

	res, err := global.DB.GetItemWithContext(ctx, input)

	if err != nil {
		return global.ErrResponse(http.StatusInternalServerError, "failed to get animal results from db"), nil
	}

	err = dynamodbattribute.UnmarshalMap(res.Item, &animal)
	if err != nil {
		return global.ErrResponse(http.StatusInternalServerError, "failed to unmarshal animal result from db"), nil
	}

	return global.Response(http.StatusOK, animal), nil
}

func main() {
	lambda.Start(find)
}
