package api

import (
	"encoding/json"
	"fmt"
	"lambda-func/database"
	"lambda-func/types"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type ApiHandler struct{
	dbStore database.UserStore
}

func NewApiHandler(dbStore database.UserStore) ApiHandler{
	return ApiHandler{
		dbStore: dbStore,
	}
}

func (api ApiHandler) RegisterUserHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error){
	var registerUser types.RegisterUser
	err := json.Unmarshal([]byte(request.Body), &registerUser)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body: "Invalid Request",
			StatusCode: http.StatusBadRequest,
		}, err
	}

	if registerUser.Username == "" || registerUser.Password == ""{
		return events.APIGatewayProxyResponse{
			Body: "Invalid Request - fields empty",
			StatusCode: http.StatusBadRequest,
		}, err
	}

	userExists, err := api.dbStore.DoesUserExists(registerUser.Username)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body: "Internal Server Error",
			StatusCode: http.StatusInternalServerError,
		}, err
	}
	if userExists {
		return events.APIGatewayProxyResponse{
			Body: "User already exists",
			StatusCode: http.StatusConflict,
		}, nil 
	}
	user, err := types.NewUser(registerUser)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body: "Internal Server Error",
			StatusCode: http.StatusInternalServerError,
		}, fmt.Errorf("could not create user %w", err)
	}
	err = api.dbStore.InsertUser(*user)
	if err != nil{
		return events.APIGatewayProxyResponse{
			Body: "Internal Server Error",
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	return events.APIGatewayProxyResponse{
		Body: "Successfully register user",
		StatusCode: http.StatusCreated,
	}, nil
}

func (api ApiHandler) LoginUserHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var loginRequest LoginRequest

	err := json.Unmarshal([]byte(request.Body), &loginRequest)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body: "Invalid request",
			StatusCode: http.StatusBadRequest,
		}, err
	}
	user, err := api.dbStore.GetUser(loginRequest.Username)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body: "Internal Server Error",
			StatusCode: http.StatusInternalServerError,
		}, err
	}
	if !types.ValidateFunction(user.PasswordHash, loginRequest.Password){
		return events.APIGatewayProxyResponse{
			Body: "Invalid credentials",
			StatusCode: http.StatusUnauthorized,
		}, nil
	}

	accessToken := types.CreateToken(user)
	successMessage := fmt.Sprintf(`{"access_token": "%s"}`, accessToken)
	return events.APIGatewayProxyResponse{
		Body: successMessage,
		StatusCode: http.StatusOK,
	}, nil
}