package main

import (
	"context"
	"fmt"
	"os"
	"time"
	"encoding/json"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/satori/go.uuid"
    
    "gopkg.in/go-playground/validator.v9"
)

type Contact struct {
	ID          string  `json:"id"`
    Name string 	`json:"name" validate:"required"`
	Email string 	`json:"email" validate:"required,email"`
	PhoneNumber string 	`json:"phone_number" validate:"required,number"`
	CreatedAt   string 	`json:"created_at"`
}

var ddb *dynamodb.DynamoDB
var validate *validator.Validate
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

func AddContact(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("AddContact")
    
	var (
        id = uuid.Must(uuid.NewV4()).String()
		tableName = aws.String(os.Getenv("CONTACTS_TABLE_NAME"))
	)
    
    validate = validator.New()
    validate.RegisterStructValidation(ContactStructLevelValidation, Contact{})

	// Initialize contact
    contact := &Contact{}
    
	/*contact := &Contact{
		ID:					id,
        Name:  "",
        Email:  "",
        PhoneNumber:  "",
		CreatedAt:			time.Now().String(),
	}*/
    
	// Parse request body
    err := json.Unmarshal([]byte(request.Body), contact)
    if err != nil{
        return events.APIGatewayProxyResponse{ // Error HTTP response
			Body: err.Error(),
			StatusCode: 500,
		}, nil
    }
    contact.ID = id
    contact.CreatedAt = time.Now().String()
    
    err = validate.Struct(contact)
    if err != nil{
        if _, ok := err.(*validator.InvalidValidationError); ok{
            fmt.Println(err)
            return events.APIGatewayProxyResponse{ // Error HTTP response
			Body: err.Error(),
			StatusCode: 500,
		      }, nil
        }
        
        for _, err := range err.(validator.ValidationErrors) {

			fmt.Println(err.Namespace())
			fmt.Println(err.Field())
			fmt.Println(err.StructNamespace()) // can differ when a custom TagNameFunc is registered or
			fmt.Println(err.StructField())     // by passing alt name to ReportError like below
			fmt.Println(err.Tag())
			fmt.Println(err.ActualTag())
			fmt.Println(err.Kind())
			fmt.Println(err.Type())
			fmt.Println(err.Value())
			fmt.Println(err.Param())
			fmt.Println()
		}
        
        return events.APIGatewayProxyResponse{ // Error HTTP response
			Body: err.Error(),
			StatusCode: 500,
		}, nil
    }


	// Write to DynamoDB
	item, _ := dynamodbattribute.MarshalMap(contact)
	input := &dynamodb.PutItemInput{
		Item: item,
		TableName: tableName,
	}
	if _, err := ddb.PutItem(input); err != nil {
		return events.APIGatewayProxyResponse{ // Error HTTP response
			Body: err.Error(),
			StatusCode: 500,
		}, nil
	} else {
		body, _ := json.Marshal(contact)
		return events.APIGatewayProxyResponse{ // Success HTTP response
			Body: string(body),
			StatusCode: 200,
		}, nil
	}
}

func ContactStructLevelValidation(sl validator.StructLevel) {

	contact := sl.Current().Interface().(Contact)

	if len(contact.Name) == 0 && len(contact.Email) == 0 && len(contact.PhoneNumber) == 0 {
		sl.ReportError(contact.Name, "Name", "name", "nameoremailorphone_number", "")
		sl.ReportError(contact.Email, "Email", "email", "nameoremailorphone_number", "")
		sl.ReportError(contact.PhoneNumber, "PhoneNumber", "phone_number", "nameoremailorphone_number", "")
	}
    
    /*if len(contact.Name) != 0 && len(contact.Email) == 0 && len(contact.PhoneNumber) == 0 {
		sl.ReportError(contact.Email, "Email", "email", "nameoremailorphone_number", "")
		sl.ReportError(contact.PhoneNumber, "PhoneNumber", "phone_number", "nameoremailorphone_number", "")
	}
    
    if len(contact.Name) == 0 && len(contact.Email) != 0 && len(contact.PhoneNumber) == 0 {
		sl.ReportError(contact.Name, "Name", "name", "nameoremailorphone_number", "")
		sl.ReportError(contact.PhoneNumber, "PhoneNumber", "phone_number", "nameoremailorphone_number", "")
	}
    
    if len(contact.Name) == 0 && len(contact.Email) == 0 && len(contact.PhoneNumber) != 0 {
		sl.ReportError(contact.Name, "Name", "name", "nameoremailorphone_number", "")
		sl.ReportError(contact.Email, "Email", "email", "nameoremailorphone_number", "")
	}*/

	// plus can to more, even with different tag than "fnameorlname"
}

func main() {
	lambda.Start(AddContact)
}