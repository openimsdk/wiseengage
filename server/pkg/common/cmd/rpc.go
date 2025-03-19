package cmd

import (
	"context"
	"fmt"
	"net"
	"reflect"
	"strconv"

	"github.com/openimsdk/tools/mw"
	"github.com/openimsdk/tools/utils/datautil"
	"github.com/openimsdk/tools/utils/network"
	"github.com/openimsdk/wiseengage/v1/pkg/common/config"
	"google.golang.org/grpc"
)

func RPCServiceRegistrar(ctx context.Context, name string, registrar func(registrar grpc.ServiceRegistrar), opt ...grpc.ServerOption) error {
	value, err := getStartContextValue(ctx)
	if err != nil {
		return err
	}
	ctx = value.RootCtx
	rpc := getSubConfig[config.RPC](reflect.ValueOf(value.Config))
	if rpc == nil {
		return fmt.Errorf("config not found rpc info")
	}
	registerIP, err := network.GetRpcRegisterIP(rpc.RegisterIP)
	if err != nil {
		return err
	}
	var rpcPort int
	if !rpc.AutoPort {
		var err error
		rpcPort, err = datautil.GetElemByIndex(rpc.Ports, value.Index)
		if err != nil {
			return err
		}
	}
	rpcListener, err := net.Listen("tcp", net.JoinHostPort(rpc.ListenIP, strconv.Itoa(rpcPort)))
	if err != nil {
		return err
	}
	defer rpcListener.Close()
	rpcPort = rpcListener.Addr().(*net.TCPAddr).Port
	rpcServer := grpc.NewServer(append(opt, mw.GrpcServer())...)
	registrar(rpcServer)
	if err := value.Client.Register(ctx, name, registerIP, rpcPort); err != nil {
		return err
	}
	serveDone := make(chan struct{})
	rpcGracefulStop := make(chan struct{})
	go func() {
		select {
		case <-serveDone:
		case <-ctx.Done():
		}
		rpcServer.GracefulStop()
		close(rpcGracefulStop)
	}()
	serveErr := rpcServer.Serve(rpcListener)
	close(serveDone)
	<-rpcGracefulStop
	return serveErr
}
