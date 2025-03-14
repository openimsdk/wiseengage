package cmd

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"reflect"
	"strconv"

	"github.com/openimsdk/tools/utils/datautil"
	"github.com/openimsdk/tools/utils/network"
	"github.com/openimsdk/wiseengage/v1/pkg/common/config"
)

func APIServer(ctx context.Context, handler http.Handler) error {
	value, err := getStartContextValue(ctx)
	if err != nil {
		return err
	}
	ctx = value.RootCtx
	api := getSubConfig[config.API](reflect.ValueOf(value.Config))
	if api == nil {
		return fmt.Errorf("config not found api info")
	}
	apiConf := api.Api
	apiPort, err := datautil.GetElemByIndex(apiConf.Ports, value.Index)
	if err != nil {
		return err
	}
	httpServer := &http.Server{
		Handler: handler,
		Addr:    net.JoinHostPort(network.GetListenIP(apiConf.ListenIP), strconv.Itoa(apiPort)),
	}
	serveDone := make(chan struct{})
	rpcGracefulStop := make(chan struct{})
	go func() {
		select {
		case <-serveDone:
		case <-ctx.Done():
		}
		_ = httpServer.Shutdown(context.Background())
		close(rpcGracefulStop)
	}()
	serveErr := httpServer.ListenAndServe()
	close(serveDone)
	<-rpcGracefulStop
	if errors.Is(serveErr, http.ErrServerClosed) {
		serveErr = nil
	}
	return serveErr
}
