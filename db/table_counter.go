package db

import (
	"context"

	"github.com/guregu/dynamo"
)

func newCounterTable(ddb *dynamo.DB, prefix string) *CounterTable {
	data := CounterTableDefinition

	t := &CounterTable{}
	t.t = NewDynamoTable(ddb, prefix, data)

	return t
}

// 账号表
// hashkey=Counter类型
type CounterTable struct {
	t *DynamoTable
}

func (t *CounterTable) createIfNotExist(ddb *dynamo.DB) error {
	return t.t.CreateTable(ddb)
}

func (t *CounterTable) Increse(ctx context.Context, key string) (int64, error) {
	// 先尝试对这个Key的值进行自增
	value := &DynamoCounter{}
	if err := t.t.UpdateStringHashKey(key).
		Add("Counter", 1).
		ValueWithContext(ctx, value); err != nil {
		return 0, err
	}
	return value.Counter, nil
}
