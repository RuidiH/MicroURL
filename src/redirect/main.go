package main

import (
	"context"
	"log"
	"os"
	"time"
	"fmt"
	"net/http"

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
	codeKeyName string
	urlKeyName  string
)

func init() {
	log.Println("init: starting Lambda lookup_redirect setup")

	tableName = os.Getenv("TABLE_NAME")
	if tableName == "" {
		log.Fatalf("init error: TABLE_NAME env must be set")
	}
	log.Printf("init: using DynamoDB table %q", tableName)

	codeKeyName = os.Getenv("CODE_KEYNAME")
	if codeKeyName == "" {
		log.Fatalf("init error: CODE_KEYNAME env must be set")
	}
	log.Printf("init: code key name %q", codeKeyName)

	urlKeyName = os.Getenv("URL_KEYNAME")
	if urlKeyName == "" {
		log.Fatalf("init error: URL_KEYNAME env must be set")
	}
	log.Printf("init: URL key name %q", urlKeyName)

	ctx := context.Background()
	awsRegion := os.Getenv("AWS_REGION")
	log.Printf("init: loading AWS SDK config (region=%q)", awsRegion)
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(awsRegion))
	if err != nil {
		log.Fatalf("init error: unable to load SDK config, %v", err)
	}

	opts := []func(*dynamodb.Options){}
	if ep := os.Getenv("DDB_ENDPOINT"); ep != "" {
		log.Printf("init: overriding DynamoDB endpoint to %q", ep)
		opts = append(opts, func(o *dynamodb.Options) {
			o.BaseEndpoint = aws.String(ep)
		})
	}
	ddbClient = dynamodb.NewFromConfig(cfg, opts...)
	log.Println("init: DynamoDB client configured for lookup_redirect")

	// wait for Lambda runtime at port 9001

    runtimeAPI := os.Getenv("AWS_LAMBDA_RUNTIME_API")
    if runtimeAPI == "" {
        // Not running under custom runtime (e.g. in AWS), so nothing to do
        return
    }

    client := &http.Client{Timeout: 5 * time.Second}
    // Try up to 5 times, 1s apart
    for i := 0; i < 5; i++ {
        resp, err := client.Get("http://" + runtimeAPI + "/2018-06-01/runtime/invocation/next")
        if err == nil {
            // API is up—close and return
            resp.Body.Close()
            fmt.Println("✅ Runtime API is ready")
            return
        }
        fmt.Printf("⏳ Runtime API not ready (%d/5): %v\n", i+1, err)
        time.Sleep(1 * time.Second)
    }

    // If you reach here, it really isn’t coming up—let Rapid handle the failure
    fmt.Println("⚠️ Runtime API never came up—continuing and will likely timeout")
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("handler: received request %+v", req)

	code := req.PathParameters["code"]
	log.Printf("handler: path parameter code = %q", code)
	if code == "" {
		log.Println("handler error: code path param missing")
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       `{"error":"code path parameter missing"}`,
		}, nil
	}

	log.Printf("handler: querying DynamoDB for code %q", code)
	out, err := ddbClient.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			codeKeyName: &types.AttributeValueMemberS{Value: code},
		},
	})
	if err != nil {
		log.Printf("handler error: DynamoDB GetItem failed: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       `{"error":"internal server error"}`,
		}, nil
	}

	if out.Item == nil {
		log.Printf("handler: code %q not found", code)
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       `{"error":"code not found"}`,
		}, nil
	}

	av, ok := out.Item[urlKeyName].(*types.AttributeValueMemberS)
	if !ok {
		log.Printf("handler error: expected string for %q but got %T", urlKeyName, out.Item[urlKeyName])
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       `{"error":"invalid data format"}`,
		}, nil
	}
	longURL := av.Value
	log.Printf("handler: found long URL %q for code %q", longURL, code)

	return events.APIGatewayProxyResponse{
		StatusCode: 302,
		Headers: map[string]string{
			"Location":                   longURL,
			"Access-Control-Allow-Origin": "*",
		},
	}, nil
}

func main() {
	log.Println("main: starting Lambda lookup_redirect")
	lambda.Start(handler)
}

