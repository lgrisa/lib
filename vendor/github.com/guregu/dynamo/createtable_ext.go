package dynamo

import "github.com/aws/aws-sdk-go/service/dynamodb"

func (ct *CreateTable) NewInput() *dynamodb.CreateTableInput {
	return ct.input()
}