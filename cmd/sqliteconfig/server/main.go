package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/lgrisa/lib/cmd/sqliteconfig/mgr"
	"github.com/pkg/errors"
	"time"
)

var (
	port   = flag.Int("port", 7787, "")
	root   = flag.String("root", "gen", "")
	prefix = flag.String("prefix", "sqlite/", "")

	region    = flag.String("region", "cn-northwest-1", "")
	bucket    = flag.String("bucket", "star-test-config", "")
	accessKey = flag.String("access_key", "", "")
	secretKey = flag.String("secret_key", "", "")
)

func main() {
	flag.Parse()

	idMapPath := fmt.Sprintf("%v/proto_id.yaml", *root)

	ak, sk := *accessKey, *secretKey
	if ak == "" {
		fmt.Println("access key not found")
		return
	}

	if sk == "" {
		fmt.Println("secret key not found")
		return
	}

	s3Client, err := initS3Client(ak, sk)
	if err != nil {
		fmt.Println("initS3Client fail", err)
		return
	}

	storage := mgr.NewStorage(*bucket, *prefix, s3Client)
	m := mgr.NewManager(*port, *root, idMapPath, storage)
	m.Start()
}

func initS3Client(accessKey, secretKey string) (*s3.Client, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(*region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				accessKey, secretKey, "")))
	if err != nil {
		return nil, errors.Wrapf(err, "LoadDefaultConfig失败")
	}

	return s3.NewFromConfig(c), nil
}
