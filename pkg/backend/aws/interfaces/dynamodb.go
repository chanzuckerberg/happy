package interfaces

import dynamolock "cirello.io/dynamolock/v2"

type DynamoDB interface {
	dynamolock.DynamoDBClient
}
