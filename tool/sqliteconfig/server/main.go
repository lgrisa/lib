package main

import (
	"flag"
	"fmt"
	"github.com/lgrisa/lib/tool/sqliteconfig/mgr"
	s4 "github.com/lgrisa/lib/tool/sqliteconfig/s3"
)

var (
	port   = flag.Int("port", 7787, "server port")
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

	s3Client, err := s4.InitS3Client(ak, sk, *region)
	if err != nil {
		fmt.Println("initS3Client fail", err)
		return
	}

	storage := s4.NewStorage(*bucket, *prefix, s3Client)
	m := mgr.NewManager(*port, *root, idMapPath, storage)
	m.Register()
}
