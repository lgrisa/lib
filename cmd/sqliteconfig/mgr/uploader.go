package mgr

import (
	"bytes"
	"context"
	oerr "errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	smithyhttp "github.com/aws/smithy-go/transport/http"
	"github.com/lgrisa/lib/cmd/sqliteconfig/mgr/pool"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"io"
)

func NewStorage(bucket, prefix string, s3Client *s3.Client) *Storage {
	return &Storage{
		bucket:   &bucket,
		prefix:   prefix,
		s3Client: s3Client,
	}
}

type Storage struct {
	bucket *string

	// 一般是data/namespace/格式. 一定要加个data前缀, 方便做replication
	prefix   string
	s3Client *s3.Client
}

func (s *Storage) Exist(ctx context.Context, key string) (bool, error) {
	// 只简单判断一下s3上是否存在, 再put
	// 不保证同时调用时能正确判断

	keyPath := s.prefix + key

	_, err := s.s3Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: s.bucket,
		Key:    aws.String(keyPath),
	})

	if err != nil {
		var responseError *smithyhttp.ResponseError
		if oerr.As(err, &responseError) {
			// 这样判断, 感觉不太稳
			if responseError.HTTPStatusCode() == 404 {
				// 对象不存在
				// 才真正开始干活
				return false, nil
			}
		}
		return false, errors.Wrap(err, "head object失败")
	} else {
		// 对象存在
		return true, nil
	}
}

func (s *Storage) Put(ctx context.Context, key string, v []byte) error {
	keyPath := s.prefix + key

	lenOfBody := int64(len(v))

	_, err := s.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        s.bucket,
		Key:           aws.String(keyPath),
		ContentLength: &lenOfBody,
		Body:          bytes.NewReader(v),
	})

	if err != nil {
		return errors.Wrapf(err, "往s3写入失败, %v", keyPath)
	}

	return nil
}

var ErrObjectNotExist = errors.Errorf("ErrObjectNotExist")

func (s *Storage) Get(ctx context.Context, key string) ([]byte, error) {
	keyPath := s.prefix + key

	object, err := s.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: s.bucket,
		Key:    aws.String(keyPath),
	})

	if err != nil {
		var notFound *types.NoSuchKey
		if oerr.As(err, &notFound) {
			// 对象不存在
			log.Debug().Msgf("读取s3.GetObject() NoSuchKey, %v", keyPath)
			return nil, ErrObjectNotExist
		}
		return nil, errors.Wrap(err, "读取s3.GetObject返回出错")
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(object.Body)

	size := int(*object.ContentLength)
	if size <= 0 {
		log.Debug().Int("size", size).Msg("s3 get object, 没有size. 只能当场通过io.ReadAll读取数据")
		data, err := io.ReadAll(object.Body)
		if err != nil {
			return nil, errors.Wrap(err, "读取s3 get返回出错")
		}

		return data, nil
	}

	log.Debug().Int("size", size).Str("key", keyPath).Msg("s3 get object, 返回包含size")
	// 有大小
	buf := pool.Pool.Alloc(size)
	if _, err := io.ReadFull(object.Body, buf); err != nil {
		buf.Free()
		return nil, errors.Wrap(err, "读取s3 get返回出错")
	}
	return buf, nil
}
