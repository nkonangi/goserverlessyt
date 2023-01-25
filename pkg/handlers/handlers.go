package handlers

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nkonangi/goserverlessyt/pkg/user"
)

var ErrorMethodNotAllowed = "method not allowed"

type ErrorBody struct {
	ErrorMsg *string `json:"error,omitempty"`
}

func GetUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI)
			(*events.APIGatewayProxyResponse,error) {
				email := req.QueryStringParameters["email"]
				if len(email) >0 {
					result, err := user.FetchUser(email, tableName, dynaClient)
					if err !=nil {
						return apiResponse(http.StatusBadRequest, ErrorBody{aws.String(err.Error()), })
					}
					return apiResponse(http.StatusOK, result)
				}


				result, err := user.FetchUsers(tableName,dynaClient)
				if err != nil {
					return apiResponse(http.StatusBadRequest,ErrorBody{ aws.String(err.Error()),})
				}
				return apiResponse(http.StatusOk, result)

}

func CreateUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI)
(*events.APIGatewayProxyResponse,error) {

	result, err := user.CreateUser(req, tableName, dynaClient)
	if err != nil {
		return apiResponse(http.StatusBadRequest,ErrorBody{
			aws.String(err.Error()),
		})
	}
	return apiResponse(http.StatusCreated, result)

}

func UpdateUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI)
(*events.APIGatewayProxyResponse,error) {

	result, err := user.UpdateUser(req, tableName, dynaClient)
	if err != nil {
		return apiResponse(http.StatusBadRequest,ErrorBody{
			aws.String(err.Error()),
		})
	}
	return apiResponse(http.StatusOk, result)

}

func DeleteUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI)
(*events.APIGatewayProxyResponse,error) {

	result, err := user.DeleteUser(req, tableName, dynaClient)
	if err != nil {
		return apiResponse(http.StatusBadRequest,ErrorBody{
			aws.String(err.Error()),
		})
	}
	return apiResponse(http.StatusOk, nil)

}

func UnhandledMethod() (*events.APIGatewayProxyResponse, error) {
	return apiResponse(http.StatusMethodNotAllowed, ErrorMethodNotAllowed)

}
