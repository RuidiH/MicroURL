package main

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var (
	ddbClient *dynamodb.Client
	tableName string
	baseURL   string
)

func init() {
	// injected by terraform
	tableName = os.Getenv("TABLE_NAME")
	if tableName == "" {
		panic("TABLE_NAME must be set")
	}

	baseURL = os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "https://short.example.com"
	}

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {
		panic("unable to load SDK config: " + err.Error())
	}

	opts := []func(*dynamodb.Options){}
	if ep := os.Getenv("DDB_ENDPOINT"); ep != "" {
		opts = append(opts, func(o *dynamodb.Options) {
			o.BaseEndpoint = aws.String(ep)
		})
	}
	ddbClient = dynamodb.NewFromConfig(cfg, opts...)
}

const (
	base62Alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	codeLength     = 7
)

func randBase62(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	n = len(base62Alphabet)
	for i := range b {
		b[i] = base62Alphabet[int(b[i])%n]
	}
	return string(b), nil
}

// generateCode tries to insert until no collision
func generateCode(ctx context.Context, longURL string) (string, error) {
	for {
		code, err := randBase62(codeLength)
		if err != nil {
			return "", fmt.Errorf("randBase62: %w", err)
		}

		// attempt a conditional write
		input := &dynamodb.PutItemInput{
			TableName: aws.String(tableName),
			Item: map[string]types.AttributeValue{
				"HashCode": &types.AttributeValueMemberS{Value: code},
				"LongURL":  &types.AttributeValueMemberS{Value: longURL},
				// "createdAt": &types.AttributeValueMemberS{
				//     Value: time.Now().UTC().Format(time.RFC3339),
				// },
			},
			ConditionExpression: aws.String("attribute_not_exists(HashCode)"),
		}
		_, err = ddbClient.PutItem(ctx, input)
		if err == nil {
			// success
			return code, nil
		}
		var ccfe *types.ConditionalCheckFailedException
		if errors.As(err, &ccfe) {
			// collision: loop and try again
			continue
		}
		// any other error: give up
		return "", fmt.Errorf("ddb PutItem: %w", err)
	}
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	longURL := request.QueryStringParameters["url"]
	if longURL == "" {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: `{"error":"url param missing"}`}, nil
	}

	code, err := generateCode(ctx, longURL)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: `{"error":"could not generate code"}`}, nil
	}

	shortURL := fmt.Sprintf("%s/%s", baseURL, code)
	respBody := fmt.Sprintf(`{"code":"%s","short_url":"%s"}`, code, shortURL)
	return events.APIGatewayProxyResponse{
		StatusCode: 201,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       respBody,
	}, nil
}

func main() {
	lambda.Start(handler)
}
