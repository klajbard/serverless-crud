package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/klajbard/serverless-crud/global"
)

// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

func findAll(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	next := event.QueryStringParameters["next"]

	result := global.AnimalsResults{
		Animals: []global.Animal{},
	}

	input := &dynamodb.ScanInput{
		TableName: aws.String(global.TableName),
		Limit:     aws.Int64(10),
	}

	if next != "" {
		input.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"id": {S: &next},
		}
	}

	res, err := global.DB.ScanWithContext(ctx, input)

	if err != nil {
		fmt.Println(err)
		return global.ErrResponse(http.StatusInternalServerError, "failed to get animal results from db"), nil
	}

	err = dynamodbattribute.UnmarshalListOfMaps(res.Items, &result.Animals)
	if err != nil {
		fmt.Println(err)
		return global.ErrResponse(http.StatusInternalServerError, "failed to unmarshal animal result from db"), nil
	}

	if len(res.LastEvaluatedKey) > 0 {
		if key, ok := res.LastEvaluatedKey["id"]; ok {
			nextKey := key.S
			result.Next = *nextKey
		}
	}

	return global.Response(http.StatusOK, result), nil
}

func main() {
	lambda.Start(findAll)
}
