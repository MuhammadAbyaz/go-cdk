package api

import (
	"fmt"
	"lambda-func/database"
	"lambda-func/types"
)

type ApiHandler struct{
	dbStore database.DynamoDBClient
}

func NewApiHandler(dbStore database.DynamoDBClient) ApiHandler{
	return ApiHandler{
		dbStore: dbStore,
	}
}

func (api ApiHandler) RegisterUserHandler(event types.RegisterUser )error{
	if event.Username == "" || event.Password == ""{
		return fmt.Errorf("request has empty parameters")
	}

	userExists, err := api.dbStore.DoesUserExists(event.Username)
	if err != nil {
		return fmt.Errorf("there was an error %w", err)
	}
	if userExists {
		return fmt.Errorf("a user with that username already exists")
	}
	err = api.dbStore.InsertUser(event)
	if err != nil{
		return fmt.Errorf("error: registering the user %w", err)
	}
	return nil
}