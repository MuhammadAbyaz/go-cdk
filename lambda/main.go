package main

import (
	"fmt"
	"lambda-func/app"
	"lambda-func/middleware"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type MyEvent struct {
	Username string `json:"username"`
}

func HandleRequest(event MyEvent) (string, error){
	if event.Username == ""{
		return "",  fmt.Errorf("username cannot be empty")
	}
	return fmt.Sprintf("Successfully called by %s", event.Username), nil
}

func ProtectedHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error){
	return events.APIGatewayProxyResponse{
		Body: "Secret Path Accessed",
		StatusCode: http.StatusOK,
	},nil
}

func main(){
	app:= app.NewApp()
	lambda.Start(func (request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		switch request.Path {
		case "/auth/register":
			return app.ApiHandler.RegisterUserHandler(request)
		case "/auth/login":
			return app.ApiHandler.LoginUserHandler(request)
		case "/protected":
			return middleware.ValidateJWTMiddleware(ProtectedHandler)(request)
		default :
			return events.APIGatewayProxyResponse{
				Body: "Not found",
				StatusCode: http.StatusNotFound,
			},nil
		}
	})
}