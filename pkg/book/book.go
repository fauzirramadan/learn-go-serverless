package book

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var(
	ErrorFailedToUnmarshalRecord = "failed to unmarshal record"
	ErrorFailedToFetchRecord = "failed to fetch record"
	ErrorInvalidBookData = "invalid book data"
	ErrorCouldNotMarshalItem = "could not marshal item"
	ErrorCouldNotDeleteItem = "could not delete item"
	ErrorCouldNotDynamoPutItem = "could not dynamo put item"
	ErrorBookAlreadyExists = "book.Book already exists"
	ErrorBookDoesNotExist = "book.Book does not exist"
)

type Book struct{
	BookName 		string	`json:"bookName"`
	Author	string 	`json:"author"`
}

func FetchBook(bookName, tableName string, dynaClient dynamodbiface.DynamoDBAPI)(*Book, error){
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email":{
				S: aws.String(bookName),
			},
		},
		TableName: aws.String(tableName),
	}
	result, err := dynaClient.GetItem(input)
	if err != nil {
		return nil, errors.New(ErrorFailedToFetchRecord)
	}
	item := new(Book)
	err = dynamodbattribute.UnmarshalMap(result.Item, item)
	if err != nil {
		return nil, errors.New(ErrorFailedToUnmarshalRecord)
	}
	return item, nil
}

func FetchBooks(tableName string, dynaClient dynamodbiface.DynamoDBAPI)(*[]Book, error){
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	result, err := dynaClient.Scan(input)
	if err!= nil {
		return nil, errors.New(ErrorFailedToFetchRecord)
	}
	item := new([]Book)
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, item)
	if err != nil {
		return nil, errors.New(ErrorFailedToUnmarshalRecord)
	}
	return item, nil
}

func CreateBook(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI)(
	*Book,
	error,
){
	var book Book

	if err := json.Unmarshal([]byte(req.Body), &book); err!=nil {
		return nil, errors.New(ErrorInvalidBookData)
	}

	currentBook, _ := FetchBook(book.BookName, tableName, dynaClient)
	if currentBook != nil && len(currentBook.BookName) != 0 {
		return nil, errors.New(ErrorBookAlreadyExists)
	}

	av, err := dynamodbattribute.MarshalMap(book)

	if err != nil {
		return nil, errors.New(ErrorCouldNotMarshalItem)
	}

	input := &dynamodb.PutItemInput{
		Item: av,
		TableName: aws.String(tableName),
	}

	_, err = dynaClient.PutItem(input)
	if err != nil {
		return nil, errors.New(ErrorCouldNotDynamoPutItem)
	}
	return &book, nil
}

func UpdateBook(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI)(
	*Book,
	error,
){
	var book Book
	if err := json.Unmarshal([]byte(req.Body), &book); err !=nil {
		return nil, errors.New(ErrorInvalidBookData)
	}

	currentUser, _ := FetchBook(book.BookName, tableName, dynaClient)
	if currentUser !=nil && len (book.BookName) == 0{
		return nil, errors.New(ErrorBookDoesNotExist)
	}

	av, err := dynamodbattribute.MarshalMap(book)
	if err != nil {
		return nil, errors.New(ErrorCouldNotMarshalItem)
	}

	input := &dynamodb.PutItemInput{
		Item: av,
		TableName: aws.String(tableName),
	}

	_, err = dynaClient.PutItem(input)
	if err!=nil {
		return nil, errors.New(ErrorCouldNotDynamoPutItem)
	}
	return &book, nil
}

func DeleteBook(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) error{
	bookName := req.QueryStringParameters["bookName"]
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email":{
				S: aws.String(bookName),
			},
		},
		TableName: aws.String(tableName),
	}
	_, err := dynaClient.DeleteItem(input)
	if err != nil {
		return errors.New(ErrorCouldNotDeleteItem)
	}

	return nil
}