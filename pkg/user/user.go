package user

import (
	"encoding/json"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/nkonangi/goserverlessyt/pkg/validations"
)

var (
	ErrorFailedToFetchRecord = "failed to fetch records"
	ErrorInvalidUserData = " Invalid user data"
	ErrorFailedToUnmarshalRecord = " failed to unmarshall record"
	ErrorInvalidEmail = "Invalid email address"
	ErrorCouldNotDeleteItem = " Could not delete this record"
	ErrorUserAlreadyExists = "User already exists"
	ErrorDoesnotExist = "user does not exist in db"
	ErrorCouldNotMarshalItem = " ErrorCouldNotMarshalItem "
	ErrorsCouldNotDynamoPutItem = "ErrorsCouldNotDynamoPutItem"


)

type User struct{
	Email string `json:"email"`
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
}

func FetchUser(email, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*User, error) {
	input := &dynamodb.GetItemInput{
		Key : map[string]*dynamodb.AttributeValue{
			"email" :{
				S: aws.String(email)
			}
		},
		TableName: aws.String(tableName)
	}
	result, err := dynaClient.GetItem(input)
	if err != nil {
		return nil, errors.New(ErrorFailedToFetchRecord)
	}

	item := new(User)
	err = dynamodbattribute.UnMarshalMap(result.Item, item)
	if err != nil {
		return nil, errors.New(ErrorFailedToUnMarshalRecord)
	}

}

func FetchUsers(tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*[]User, error) {

	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName)
	}

	result, err := dynaClient.Scan(input)
	if err != nil {
		return nil, errors.New(ErrorFailedToFetchRecord)
	}

	item := new({}User)
	err = dynamodbattribute.UnMarshalMap(result.Items, item)
	return item, nil

}

func CreateUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI)(*User,error) {

	var u User

	err := json.Unmarshal([]byte(req.Body), &u)

	if err != nil {
		return nil, errors.New(ErrorInvalidUserData)
	}

	if ! validations.IsEmailValid(u.Email){
		return nil, errors.New("Invalid Email")
	}
	// check if the user existis

	currentUser , _ := FetchUser(u.Email,tableName,dynaClient)
	if currentUser !=nil && len(currentUser.Email)!=0 {
		return nil, errors.New(ErrorUserAlreadyExists)
	}
	
	av, err := dynamodbattribute.marshalMap(u)

	if err !=nil {
		return nil, errors.New(ErrorCouldNotMarshalItem)
	}

	input := &dynamodb.PutItemInput{
		ITem : av,
		TableName : aws.String(tableName)
	}

	_, err:= dynaClient.PutItem(input)

	if err != nil{
		return nil, errors.New(ErrorsCouldNotDynamoPutItem)
	}

	return &u, nil

}

func UpdateUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI)(*User,error) {
	var u User

	err := json.Unmarshal([]byte(req.Body), &u)

	if err != nil {
		return nil, errors.New(ErrorInvalidUserData)
	}

	currentUser , _ := FetchUser(u.Email,tableName,dynaClient)
	if currentUser !=nil && len(currentUser.Email)==0 {
		return nil, errors.New(ErrorDoesnotExist)
	}

	av, err := dynamodbattributes.MarshalMap(u)
	if err != nil {
		return nil, errors.New(ErrorCouldNotMarshalItem)
	}

	input := &dynamodb.PutItemInput{
		Item : av,
		TableName : aws.String(tableName)
	}

	_,err := dynaClient.PutItem(item)

	if err != nil {
		return nil, errors.New(ErrorsCouldNotDynamoPutItem)
	}
}

func DeleteUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) error {

	email := req.QueryStringParameters["email"]
	input := & dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email" :{
				S:aws.String(email)
			}
		}
		Tablename: aws.String(tableName)
	}
	_, err := dynaClient.DeleteItem(input)
	if err !=nil {
		return errors.New(ErrorCouldNotDeleteItem)
	}
	return nil
}
