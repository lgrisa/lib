package db

import (
	"github.com/lgrisa/lib/dynamo/db/dbdef"
)

var CounterTableDefinition = &dbdef.TableDefinition{
	TableName: "Counter",
	HashKey:   "Id",
	NewEntity: func() interface{} {
		return &DynamoCounter{}
	},
}

// DynamoCounter 自增表
type DynamoCounter struct {

	// Id
	Id string `dynamo:"Id,hash"`

	Counter int64 `dynamo:"Counter"`
}
