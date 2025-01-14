package db

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/guregu/dynamo"
	"github.com/lgrisa/lib/dynamo/db/dbdef"
	"github.com/lgrisa/lib/utils/logutil"
	"github.com/pkg/errors"
	"strings"
)

func NewDynamoTable(ddb *dynamo.DB, prefix string, definition *dbdef.TableDefinition) *DynamoTable {
	table := ddb.Table(definition.NewTableName(prefix))
	return newDynamoTable(&table, definition.HashKey, definition.RangeKey, definition.TtlKey, definition.NewEntity)
}

//func newHashDynamoTable(t *dynamo.Table, hashKeyName string) *DynamoTable {
//	return newDynamoTable(t, hashKeyName, "", "")
//}

func newDynamoTable(t *dynamo.Table, hashKeyName, rangeKeyName, ttlKeyName string, newEntity func() interface{}) *DynamoTable {
	table := &DynamoTable{
		Table:     t,
		newEntity: newEntity,
		HashKey: &DynamoKey{
			KeyName: hashKeyName,
			//KeyType: hashKeyType,
		},
	}

	table.HashKeyAttributeExist = "attribute_exists(" + hashKeyName + ")"
	table.HashKeyAttributeNotExist = "attribute_not_exists(" + hashKeyName + ")"

	if rangeKeyName != "" {
		table.RangeKey = &DynamoKey{
			KeyName: rangeKeyName,
			//KeyType: rangeKeyType,
		}
	}

	if ttlKeyName != "" {
		table.TtlKey = &DynamoKey{
			KeyName: ttlKeyName,
			//KeyType: ttlKeyType,
		}
	}

	return table
}

type DynamoTable struct {
	*dynamo.Table

	newEntity func() interface{}

	HashKey  *DynamoKey
	RangeKey *DynamoKey
	TtlKey   *DynamoKey

	HashKeyAttributeExist    string
	HashKeyAttributeNotExist string
}

func (t *DynamoTable) CreateTable(ddb *dynamo.DB) error {
	return t.CreateTableWithDefinition(ddb, t.newEntity())
}

func (t *DynamoTable) CreateTableWithDefinition(ddb *dynamo.DB, from interface{}) error {

	tableExist := false

	tableName := t.Name()
	err := ddb.CreateTable(tableName, from).OnDemand(true).Run()
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == dynamodb.ErrCodeResourceInUseException {
				logutil.LogInfoF("创建表: %v, 表已经存在，跳过", tableName)
				tableExist = true
				goto ttl
			}
		}
		return errors.Wrapf(err, tableName)
	}
	logutil.LogInfoF("创建表: %v 成功", tableName)

ttl:

	if t.TtlKey != nil {
		if ttl, err := t.DescribeTTL().Run(); err != nil {
			return errors.Wrapf(err, tableName)
		} else {

			if ttl.Attribute != "" && ttl.Attribute != t.TtlKey.KeyName {
				return errors.Errorf("创建表: %v, ttl属性不匹配，期望: %v, 实际: %v", tableName, t.TtlKey.KeyName, ttl.Attribute)
			}

			switch ttl.Status {
			case dynamo.TTLEnabled:
				logutil.LogInfoF("创建表: %v, ttl已经启用，跳过", tableName)
			case dynamo.TTLDisabled:
				if err := t.UpdateTTL(t.TtlKey.KeyName, true).Run(); err != nil {
					return errors.Wrapf(err, "创建表: %v, 启用ttl出错", tableName)
				} else {
					logutil.LogInfoF("创建表: %v, 启用ttl成功", tableName)
				}
			case dynamo.TTLEnabling:
				logutil.LogInfoF("创建表: %v, ttl正在启用，跳过", tableName)
			case dynamo.TTLDisabling:
				logutil.LogInfoF("创建表: %v, ttl正在禁用，跳过", tableName)
			default:
				return errors.Errorf("创建表: %v, ttl状态未知: %v", tableName, ttl.Status)
			}
		}
	}

	if tableExist {
		createTableInput := ddb.CreateTable(tableName, from).OnDemand(true).NewInput()

		if len(createTableInput.GlobalSecondaryIndexes) > 0 {

			existIndexMap := map[string]struct{}{}
			if desc, err := t.Describe().Run(); err != nil {
				return errors.Wrapf(err, "创建索引[%v].[%v]失败", tableName, *createTableInput)
			} else {
				for _, gsi := range desc.GSI {
					existIndexMap[gsi.Name] = struct{}{}
				}
			}

			for _, gsi := range createTableInput.GlobalSecondaryIndexes {

				if _, ok := existIndexMap[*gsi.IndexName]; ok {
					logutil.LogInfoF("创建表[%v]的索引[%v]，索引已经存在，跳过", tableName, *gsi.IndexName)
					continue
				}

				if resp, err := ddb.Client().UpdateTable(&dynamodb.UpdateTableInput{
					TableName:            aws.String(tableName),
					AttributeDefinitions: createTableInput.AttributeDefinitions,
					BillingMode:          createTableInput.BillingMode,
					GlobalSecondaryIndexUpdates: []*dynamodb.GlobalSecondaryIndexUpdate{
						{
							Create: &dynamodb.CreateGlobalSecondaryIndexAction{
								IndexName:             gsi.IndexName,
								KeySchema:             gsi.KeySchema,
								Projection:            gsi.Projection,
								ProvisionedThroughput: gsi.ProvisionedThroughput,
							},
						},
					},
				}); err != nil {
					if aErr, ok := err.(awserr.Error); ok {
						if strings.Contains(strings.ToLower(aErr.Message()), "exists") {
							logutil.LogInfoF("创建表[%v]的索引[%v]，索引已经存在，跳过, err: %v", tableName, *gsi.IndexName, err)
							continue
						}
					}
					return errors.Wrapf(err, "创建索引[%v].[%v]失败", tableName, *gsi.IndexName)
				} else {
					logutil.LogInfoF("创建表[%v]的索引[%v]成功, resp: %+v", tableName, *gsi.IndexName, resp)
				}
			}
		}
	}

	return nil
}

func (t *DynamoTable) LoadStringKey(ctx aws.Context, hashKey, rangeKey string, out interface{}) error {
	return t.Table.Get(t.HashKey.KeyName, hashKey).
		Range(t.RangeKey.KeyName, dynamo.Equal, rangeKey).
		OneWithContext(ctx, out)
}

func (t *DynamoTable) LoadStringHashKey(ctx aws.Context, hashKey string, out interface{}) error {
	return t.Table.Get(t.HashKey.KeyName, hashKey).
		OneWithContext(ctx, out)
}

func (t *DynamoTable) LoadInt64HashKey(ctx aws.Context, hashKey int64, out interface{}) error {
	return t.Table.Get(t.HashKey.KeyName, hashKey).
		OneWithContext(ctx, out)
}

func (t *DynamoTable) UpdateStringHashKey(hashKey string) *dynamo.Update {
	return t.Table.Update(t.HashKey.KeyName, hashKey)
}

func (t *DynamoTable) DeleteStringKey(hashKey, rangeKey string) *dynamo.Delete {
	return t.Table.Delete(t.HashKey.KeyName, hashKey).Range(t.RangeKey.KeyName, rangeKey)
}

func (t *DynamoTable) DeleteInt64HashKey(hashKey int64) *dynamo.Delete {
	return t.Table.Delete(t.HashKey.KeyName, hashKey)
}

func (t *DynamoTable) DeleteStringHashKey(hashKey string) *dynamo.Delete {
	return t.Table.Delete(t.HashKey.KeyName, hashKey)
}

func (t *DynamoTable) BatchWrite() *dynamo.BatchWrite {
	if t.RangeKey == nil {
		return t.Table.Batch(t.HashKey.KeyName).Write()
	}

	return t.Table.Batch(t.HashKey.KeyName, t.RangeKey.KeyName).Write()
}

type DynamoKey struct {
	KeyName string
}
