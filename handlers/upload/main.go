package main

import (
	"context"
	"fmt"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/klajbard/serverless-crud/global"
)

// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

func upload(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := event.QueryStringParameters["id"]

	if id == "" {
		return global.ErrResponse(http.StatusBadRequest, "no item specified id to find"), nil
	}

	var boundary string
	if _, params, err := mime.ParseMediaType(event.Headers["content-type"]); err == nil {
		boundary = params["boundary"]
	}

	reader := multipart.NewReader(strings.NewReader(event.Body), boundary)
	part, err := reader.NextPart()

	if err != nil {
		fmt.Println(err)
		return global.ErrResponse(http.StatusBadRequest, "failed to parse form data"), nil
	}

	contentType := "image/jpeg"
	if strings.HasSuffix(part.FileName(), ".png") {
		contentType = "image/png"
	}

	uploader := s3manager.NewUploader(global.Session)
	output, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(os.Getenv("ANIMALS_BUCKET_NAME")),
		Key:         aws.String(part.FileName()),
		Body:        part,
		ContentType: aws.String(contentType),
	})

	fmt.Println(output.Location)

	if err != nil {
		fmt.Println(err)
		return global.ErrResponse(http.StatusInternalServerError, "Failed to upload image to S3"), nil
	}

	return global.Response(http.StatusOK, global.SimpleResponse{Response: fmt.Sprintf("successfully uploaded: %s", part.FileName())}), nil
}

func main() {
	lambda.Start(upload)
}
