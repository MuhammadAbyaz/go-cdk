package database

import (
	"fmt"
	"lambda-func/types"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const (
	TABLE_NAME="userTable"
)

type UserStore interface{
	DoesUserExists(username string)(bool,error)
	InsertUser(user types.User) error
	GetUser(username string) (types.User, error)
}
type DynamoDBClient struct {
	databaseStore *dynamodb.DynamoDB
}


func NewDynamoDBClient () DynamoDBClient{
	dbSession := session.Must(session.NewSession())
	db:= dynamodb.New(dbSession)
	return DynamoDBClient{
		databaseStore: db,
	}
}
func(u DynamoDBClient) DoesUserExists(username string)(bool, error){
	result, err := u.databaseStore.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(TABLE_NAME),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(username),
			},
		},
	})
	if err != nil {
		return false, err
	}
	if result.Item == nil {
		return false, nil
	}
	return true, nil
}

func (u DynamoDBClient) GetUser(username string) (types.User, error){
	var user types.User
	res, err := u.databaseStore.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(TABLE_NAME),
		Key: map[string]*dynamodb.AttributeValue{
			"username" : {
				S: aws.String(username),
			},
		},
	})
	if err != nil {
		return user, err
	}
	if res.Item == nil {
		return user, fmt.Errorf("user not found")
	}
	err = dynamodbattribute.UnmarshalMap(res.Item, &user)
	if err != nil{
		return user, err
	}
	return user, nil
}

func (u DynamoDBClient) InsertUser(user types.User) error{
	item:= &dynamodb.PutItemInput{
		TableName: aws.String(TABLE_NAME),
		Item: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(user.Username),
			},
			"password": {
				S: aws.String(user.PasswordHash),
			},
		},
	}
	_, err := u.databaseStore.PutItem(item)
	if err != nil {
		return err
	}
	return nil
}