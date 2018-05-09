package main

import (
	"errors"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	dynamo *dynamodb.DynamoDB

	region = os.Getenv("AWS_REGION")

	tableName = aws.String(os.Getenv("PRODUCTS_TABLE_NAME"))

	// ErrNameNotProvided is thrown when a name is not provided
	ErrNameNotProvided = errors.New("no name was provided in the HTTP body")
)

func init() {
	log.Println("Calling init")

	log.Printf("Region: %s, Table %s\n", region, *tableName)

	// Use aws sdk to connect to dynamoDB
	if session, err := session.NewSession(&aws.Config{Region: &region}); err != nil {
		log.Printf("Failed to connect to AWS NewSession: %s\n", err.Error())
	} else {
		dynamo = dynamodb.New(session) // Create DynamoDB client
	}

}

// HandlerHello is your Lambda function handler
// It uses Amazon API Gateway request/responses provided by the aws-lambda-go/events package,
// However you could use other event sources (S3, Kinesis etc), or JSON-decoded primitive types such as 'string'.
func HandlerHello(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// stdout and stderr are sent to AWS CloudWatch Logs
	log.Printf("Processing Lambda request %s\n", request.RequestContext.RequestID)

	log.Print("Hello from LambdaHandler?")

	log.Print("Request = " + request.RequestContext.RequestID)

	// Read from DynamoDB
	input := &dynamodb.ScanInput{
		TableName: tableName,
	}

	if result, err := dynamo.Scan(input); err != nil {
		log.Printf("Failed to Scan: %s\n", err.Error())
	} else {
		for _, i := range result.Items {
			log.Printf("%+v\n", i)
		}
	}

	// If no name is provided in the HTTP request body, throw an error
	if len(request.Body) < 1 {
		return events.APIGatewayProxyResponse{}, ErrNameNotProvided
	}

	response := events.APIGatewayProxyResponse{
		Body:       "Helloo " + request.Body,
		StatusCode: 200,
	}

	log.Printf("%+v\n",response)

	return response, nil

}

func main() {
	log.Print("Hello from Lambda")
	lambda.Start(HandlerHello)
}
