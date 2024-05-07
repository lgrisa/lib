package rpc

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"strconv"
	"sync"
)

var (
	cachedServer *RpcServer
	cachedErr    error
	createOnce   sync.Once
)

// 接收来自其他服的grpc连接
type RpcServer struct {
	rpcServer           *grpc.Server
	listener            net.Listener
	rpcServerClosedChan chan struct{}
	closeOnce           sync.Once
	startOnce           sync.Once
}

func GetOrNewRpcServer() (*RpcServer, error) {
	createOnce.Do(func() {
		port := startconfig.StartConfig.RpcPort
		listener, err := net.Listen("tcp", ":"+strconv.Itoa(int(port)))
		if err != nil {
			cachedErr = errors.Wrapf(err, "监听rpc端口失败: %d", port)
			return
		}

		rpcServer := grpc.NewServer()

		result := &RpcServer{
			rpcServer:           rpcServer,
			listener:            listener,
			rpcServerClosedChan: make(chan struct{}),
		}

		cachedServer, cachedErr = result, nil
		return
	})

	return cachedServer, cachedErr
}

func (r *RpcServer) Close() {
	r.closeOnce.Do(func() {
		close(r.rpcServerClosedChan)
		logrus.Info("rpc service 退出")
		r.listener.Close()
		//r.rpcServer.Close()
	})
}

func (r *RpcServer) ClosedChan() <-chan struct{} {
	return r.rpcServerClosedChan
}

func (r *RpcServer) RawRpcServer() *grpc.Server {
	return r.rpcServer
}

func (r *RpcServer) Start() {
	r.startOnce.Do(func() {

		//go call.CatchLoopPanic("rpc service", func() {
		//	defer r.Close()
		//	if err := r.rpcServer.Serve("memu", "local"); err != nil {
		//		logrus.WithError(err).Errorf("rpc service 监听memu失败")
		//	}
		//})

		go call.CatchLoopPanic("rpc service", func() {
			defer r.Close()
			logrus.Infof("rpc service 启动: %s", r.listener.Addr().String())

			if err := r.rpcServer.Serve(r.listener); err != nil {
				logrus.WithError(err).Errorf("rpc服务器退出")
			}
		})
	})
}
