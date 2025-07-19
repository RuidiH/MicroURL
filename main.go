package main

import (
	"crypto/sha256"
	"math/big"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func init() {
	// do any initialization here
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
	for i, j := 0, len(chars) - 1; i < j; i, j = i + 1, j - 1 {
		chars[i], chars[j] = chars[j], chars[i]
	}
	return string(chars)
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	longURL := request.QueryStringParameters["url"]
	if longURL == "" {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: `{"error":"url param missing"}`}, nil
	}

	hash := sha256.Sum256([]byte(longURL))	
	code := base62Encode(hash[:5])
	short := "https://short.ly/" + code
	body := `{"code":"` + code + `","short_url":"` + short + `"}`

	return events.APIGatewayProxyResponse{
		StatusCode:	201,
		Headers: map[string]string {
			"Content-Type"	: "application/json",
			"Location"		: short,
			"Access-Control-Allow-Origin": "*",
		},
		Body:	body,
	}, nil
}

func main() {
	lambda.Start(handler)
}
