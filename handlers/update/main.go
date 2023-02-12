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
	"github.com/klajbard/serverless-crud/global"
)

// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

func update(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := event.QueryStringParameters["id"]

	if id == "" {
		return global.ErrResponse(http.StatusBadRequest, "no item specified to update"), nil
	}

	if strings.TrimSpace(event.Body) == "" {
		return global.ErrResponse(http.StatusBadRequest, "empty request body"), nil
	}

	animal := global.Animal{}

	if err := json.Unmarshal([]byte(event.Body), &animal); err != nil {
		fmt.Println(err)
		return global.ErrResponse(http.StatusBadRequest, "failed to parse request body"), nil
	}

	_, err := dynamodbattribute.MarshalMap(&animal)
	if err != nil {
		fmt.Println(err)
		return global.ErrResponse(http.StatusInternalServerError, "unable to marshal animal data"), nil
	}

	var updateExpressions []string
	updateExpressions = append(updateExpressions, "set #status=:status")
	if animal.Name != "" {
		updateExpressions = append(updateExpressions, "#name=:name")
	}
	if animal.Description != "" {
		updateExpressions = append(updateExpressions, "description=:description")
	}
	if animal.Avatar != "" {
		updateExpressions = append(updateExpressions, "avatar=:avatar")
	}
	if animal.Breed != "" {
		updateExpressions = append(updateExpressions, "breed=:breed")
	}
	if animal.Birth != 0 {
		updateExpressions = append(updateExpressions, "birth=:birth")
	}

	updateData, err := dynamodbattribute.MarshalMap(global.AnimalUpdate{
		Name:        animal.Name,
		Breed:       animal.Breed,
		Description: animal.Description,
		Status:      animal.Status,
		Avatar:      animal.Avatar,
		Birth:       animal.Birth,
	})

	if err != nil {
		fmt.Println(err)
		return global.ErrResponse(http.StatusInternalServerError, "failed to update animal results from db"), nil
	}

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(global.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {S: &id},
		},
		ConditionExpression:       aws.String("attribute_exists(id)"),
		ExpressionAttributeValues: updateData,
		ExpressionAttributeNames: map[string]*string{
			"#name":   aws.String("name"),
			"#status": aws.String("status"),
		},
		UpdateExpression: aws.String(strings.Join(updateExpressions[:], ", ")),
	}

	_, err = global.DB.UpdateItemWithContext(ctx, input)

	if err != nil {
		fmt.Println(err)
		return global.ErrResponse(http.StatusInternalServerError, "failed to update animal results from db"), nil
	}

	return global.Response(http.StatusOK, global.SimpleResponse{Response: "successfully updated"}), nil
}

func main() {
	lambda.Start(update)
}
