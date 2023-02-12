package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/klajbard/serverless-crud/global"
)

// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

func remove(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := event.QueryStringParameters["id"]

	if id == "" {
		return global.ErrResponse(http.StatusBadRequest, "no item id specified to remove"), nil
	}

	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(global.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {S: &id},
		},
	}

	_, err := global.DB.DeleteItemWithContext(ctx, input)

	if err != nil {
		fmt.Println(err)
		return global.ErrResponse(http.StatusInternalServerError, "failed to remove animal results from db"), nil
	}

	return global.Response(http.StatusOK, global.SimpleResponse{Response: "successfully removed"}), nil
}

func main() {
	lambda.Start(remove)
}
