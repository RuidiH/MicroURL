package main

import (
	"context"
	"crypto/sha256"
	"math/big"
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
)

func init() {
	// injected by terraform
	tableName = os.Getenv("TABLE_NAME")
	if tableName == "" {
		panic("TABLE_NAME must be set")
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

const base62Alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func base62Encode(b []byte) string {
	num := new(big.Int).SetBytes(b)
	zero := big.NewInt(0)
	base := big.NewInt(62)

	var chars []byte // remainders
	for num.Cmp(zero) > 0 {
		rem := new(big.Int)
		num.DivMod(num, base, rem)
		chars = append(chars, base62Alphabet[rem.Int64()])
	}

	// if hash is all zeros
	if len(chars) == 0 {
		return "0"
	}

	// swap chars values
	for i, j := 0, len(chars)-1; i < j; i, j = i+1, j-1 {
		chars[i], chars[j] = chars[j], chars[i]
	}
	return string(chars)
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	longURL := request.QueryStringParameters["url"]
	if longURL == "" {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: `{"error":"url param missing"}`}, nil
	}

	hash := sha256.Sum256([]byte(longURL))
	code := base62Encode(hash[:5])
	short := "https://short.ly/" + code
	body := `{"code":"` + code + `","short_url":"` + short + `"}`

	item := map[string]types.AttributeValue{
		"HashCode":     &types.AttributeValueMemberS{Value: code},
		"LongURL": &types.AttributeValueMemberS{Value: longURL},
	}

	_, err := ddbClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	})
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: `{"error":"write_failed"}`}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 201,
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Location":                    short,
			"Access-Control-Allow-Origin": "*",
		},
		Body: body,
	}, nil
}

func main() {
	lambda.Start(handler)
}
