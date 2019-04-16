package main

import (
	"context"
	"fmt"
	"encoding/json"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Contact struct {
	ID          string  `json:"id"`
	Name string 	`json:"name"`
	Email string 	`json:"email"`
	PhoneNumber string 	`json:"phone_number"`
	CreatedAt   string 	`json:"created_at"`
}

type ListContactsResponse struct {
	Contacts		[]Contact  `json:"contacts"`
}

var ddb *dynamodb.DynamoDB
func init() {
	region := os.Getenv("AWS_REGION")
	if session, err := session.NewSession(&aws.Config{ // Use aws sdk to connect to dynamoDB
		Region: &region,
	}); err != nil {
		fmt.Println(fmt.Sprintf("Failed to connect to AWS: %s", err.Error()))
	} else {
		ddb = dynamodb.New(session) // Create DynamoDB client
	}
}

func ListContacts(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("ListContacts")

	var (
		tableName = aws.String(os.Getenv("CONTACTS_TABLE_NAME"))
	)

	// Read from DynamoDB
	input := &dynamodb.ScanInput{
		TableName: tableName,
	}
	result, _ := ddb.Scan(input)

	// Construct contacts from response
	var contacts []Contact
	for _, i := range result.Items {
		contact := Contact{}
		if err := dynamodbattribute.UnmarshalMap(i, &contact); err != nil {
			fmt.Println("Failed to unmarshal")
			fmt.Println(err)
		}
		contacts = append(contacts, contact)
	}

	// Success HTTP response
	body, _ := json.Marshal(&ListContactsResponse{
		Contacts: contacts,
	})
	return events.APIGatewayProxyResponse{
		Body: string(body),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(ListContacts)
}