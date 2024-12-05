package s3

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/pkg/errors"
	"time"
)

func InitS3Client(accessKey, secretKey, region string) (*s3.Client, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				accessKey, secretKey, "")))
	if err != nil {
		return nil, errors.Wrapf(err, "LoadDefaultConfig失败")
	}

	return s3.NewFromConfig(c), nil
}
