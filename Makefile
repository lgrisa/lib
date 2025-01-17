proto:
	cd sbaseserverproto && protoc --gofast_out=. *.proto

wechat_ext:
	cp vendor_ext/wechat_ext.gobak vendor/github.com/go-pay/gopay/wechat/v3/wechat_ext.go

goMod:
	go mod tidy
	go mod vendor
	cp vendor_ext/createtable_ext.gobak vendor/github.com/guregu/dynamo/createtable_ext.go

reset:
	git reset --hard HEAD
	git clean -f -d
	git pull

ctest:
	go test ./...

#go build -gcflags -m main.go 