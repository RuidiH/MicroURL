package main

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var (
	ddbClient   *dynamodb.Client
	tableName   string
	baseURL     string
	codeKeyName string
	urlKeyName  string
)

func init() {
	log.Println("init: starting Lambda setup")
	// injected by terraform
	tableName = os.Getenv("TABLE_NAME")
	if tableName == "" {
		log.Fatalf("init error: TABLE_NAME env not found")
	}
	log.Printf("init: using DynamoDB table %q", tableName)

	baseURL = os.Getenv("BASE_URL")
	if baseURL == "" {
		log.Printf("init warning: BASE_URL not set, defaulting to short.example.com")
		baseURL = "https://short.example.com"
	} else {
		log.Printf("init: using BASE_URL %q", baseURL)
	}

	codeKeyName = os.Getenv("CODE_KEYNAME")
	if codeKeyName == "" {
		log.Fatalf("init error: CODE_KEYNAME env not found")
	}
	log.Printf("init: using code key name %q", codeKeyName)

	urlKeyName = os.Getenv("URL_KEYNAME")
	if urlKeyName == "" {
		log.Fatalf("init error: URL_KEYNAME env not found")
	}
	log.Printf("init: using URL key name %q", urlKeyName)

	ctx := context.Background()
	awsRegion := os.Getenv("AWS_REGION")
	log.Printf("init: loading AWS SDK config (region=%q)", awsRegion)
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(awsRegion))
	if err != nil {
		log.Fatalf("init error: unable to load AWS SDK config: %v", err)
	}

	var opts []func(*dynamodb.Options)
	if ep := os.Getenv("DDB_ENDPOINT"); ep != "" {
		log.Printf("init: overriding DynamoDB endpoint to %q", ep)
		opts = append(opts, func(o *dynamodb.Options) {
			o.BaseEndpoint = aws.String(ep)
		})
	}
	ddbClient = dynamodb.NewFromConfig(cfg, opts...)
	log.Println("init: DynamoDB client configured")
}

const (
	base62Alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	codeLength     = 7
)

var payload struct {
	LongURL string `json:"url"`
}

func randBase62(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		log.Printf("randBase62 error: %v", err)
		return "", err
	}
	alphabetLen := len(base62Alphabet)
	for i := range b {
		b[i] = base62Alphabet[int(b[i])%alphabetLen]
	}
	return string(b), nil
}

// generateCode tries to insert until no collision
func generateCode(ctx context.Context, longURL string) (string, error) {
	attempt := 0
	for {
		attempt++
		code, err := randBase62(codeLength)
		if err != nil {
			log.Printf("generateCode: attempt %d: randBase62 error: %v", attempt, err)
			return "", fmt.Errorf("randBase62: %w", err)
		}
		log.Printf("generateCode: attempt %d: trying code %q", attempt, code)

		input := &dynamodb.PutItemInput{
			TableName: aws.String(tableName),
			Item: map[string]types.AttributeValue{
				codeKeyName: &types.AttributeValueMemberS{Value: code},
				urlKeyName:  &types.AttributeValueMemberS{Value: longURL},
			},
			ConditionExpression: aws.String("attribute_not_exists(" + codeKeyName + ")"),
		}
		_, err = ddbClient.PutItem(ctx, input)
		if err == nil {
			log.Printf("generateCode: success on attempt %d with code %q", attempt, code)
			return code, nil
		}
		var ccfe *types.ConditionalCheckFailedException
		if errors.As(err, &ccfe) {
			log.Printf("generateCode: collision on code %q, retrying", code)
			continue
		}
		log.Printf("generateCode: unexpected PutItem error: %v", err)
		return "", fmt.Errorf("ddb PutItem: %w", err)
	}
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("handler: received request %+v", request)

	if err := json.Unmarshal([]byte(request.Body), &payload); err != nil || payload.LongURL == "" {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: `{"error":"url required"}`}, nil
	}

	longURL := payload.LongURL
	log.Printf("handler: long_url param = %q", longURL)
	if longURL == "" {
		log.Printf("handler error: url param missing")
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: `{"error":"url param missing"}`}, nil
	}

	code, err := generateCode(ctx, longURL)
	if err != nil {
		log.Printf("handler error: code generation failed: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: `{"error":"could not generate code"}`}, nil
	}

	shortURL := fmt.Sprintf("%s/%s", baseURL, code)
	respBody := fmt.Sprintf(`{"code":"%s","short_url":"%s"}`, code, shortURL)
	log.Printf("handler: generated short_url %q for long_url %q", shortURL, longURL)

	return events.APIGatewayProxyResponse{
		StatusCode: 201,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       respBody,
	}, nil
}

func main() {
	log.Println("main: starting Lambda create_url")
	lambda.Start(handler)
}
