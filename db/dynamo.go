package db

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"github.com/packer/config"
	"github.com/packer/dbv"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type DynamoDB struct {
	Db      *dynamo.DB
	dbTable dynamo.Table

	userIdCounter *CounterTable

	accountTable *DynamoTable
}

func NewDynamoClient() (*DynamoDB, error) {
	c := config.StartConfig.Aws
	if c.AwsRegion == "" {
		return nil, errors.New("aws_region is empty")
	}

	awsConfig := &aws.Config{}
	awsConfig.Region = aws.String(c.AwsRegion)

	if c.AccessKey != "" && c.SecretKey != "" {
		awsConfig.Credentials = credentials.NewStaticCredentials(c.AccessKey, c.SecretKey, "")
	} else if c.Profile != "" {
		awsConfig.Credentials = credentials.NewSharedCredentials("", c.Profile)
	} else {
		awsConfig.Credentials = credentials.NewEnvCredentials()
	}

	if c.DynamoEndpoint != "" {
		awsConfig.Endpoint = aws.String(c.DynamoEndpoint)
	}

	sess := session.Must(session.NewSession(awsConfig))
	db := dynamo.New(sess)

	var tableArray []*DynamoTable
	addTable := func(t *DynamoTable) *DynamoTable {
		tableArray = append(tableArray, t)
		return t
	}

	accountTable := addTable(NewDynamoTable(db, c.DynamoTablePrefix, dbv.AccountTableDefinition))

	userIdCounter := newCounterTable(db, c.DynamoTablePrefix)

	if c.CreateTableAnyway {
		logrus.Infof("尝试创建db表")

		for _, t := range tableArray {
			if err := t.CreateTable(db); err != nil {
				return nil, errors.Wrapf(err, "创建表失败")
			}
		}

		if err := userIdCounter.t.CreateTable(db); err != nil {
			return nil, errors.Wrapf(err, "创建表失败")
		}
	}

	return &DynamoDB{
		Db:            db,
		userIdCounter: userIdCounter,
		accountTable:  accountTable,
	}, nil
}

func (d *DynamoDB) AccountTable() *DynamoTable {
	return d.accountTable
}