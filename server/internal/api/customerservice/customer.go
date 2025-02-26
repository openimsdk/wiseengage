package customerservice

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/network"
	"github.com/openimsdk/tools/utils/runtimeenv"
	"github.com/openimsdk/wiseengage/v1/pkg/common/config"
	"github.com/openimsdk/wiseengage/v1/pkg/protocol/customerservice"
	"google.golang.org/grpc"
)

type Config struct {
	RpcService config.RpcService
	API        config.API
}

func Start(ctx context.Context, config *Config, client discovery.Conn, service grpc.ServiceRegistrar) error {
	router, err := newGinRouter(ctx, client, config)
	if err != nil {
		return err
	}
	apiCtx, apiCancel := context.WithCancelCause(context.Background())
	done := make(chan struct{})
	go func() {
		httpServer := &http.Server{
			Handler: router,
			Addr:    net.JoinHostPort(network.GetListenIP(config.API.Api.ListenIP), strconv.Itoa(config.API.Api.Port)),
		}
		go func() {
			defer close(done)
			select {
			case <-ctx.Done():
				apiCancel(fmt.Errorf("recv ctx %w", context.Cause(ctx)))
			case <-apiCtx.Done():
			}
			log.ZDebug(ctx, "api server is shutting down")
			if err := httpServer.Shutdown(context.Background()); err != nil {
				log.ZWarn(ctx, "api server shutdown err", err)
			}
		}()
		log.CInfo(ctx, "api server is init", "runtimeEnv", runtimeenv.RuntimeEnvironment(), "address", httpServer.Addr, "apiPort", config.API.Api.Port)
		err := httpServer.ListenAndServe()
		if err == nil {
			err = errors.New("api done")
		}
		apiCancel(err)
	}()

	<-apiCtx.Done()
	exitCause := context.Cause(ctx)
	log.ZWarn(ctx, "api server exit", exitCause)
	timer := time.NewTimer(time.Second * 15)
	defer timer.Stop()
	select {
	case <-timer.C:
		log.ZWarn(ctx, "api server graceful stop timeout", nil)
	case <-done:
		log.ZDebug(ctx, "api server graceful stop done")
	}
	return exitCause
}

func newGinRouter(ctx context.Context, client discovery.Conn, cfg *Config) (*gin.Engine, error) {
	wiseEngageConn, err := client.GetConn(ctx, cfg.RpcService.WiseEngage)
	if err != nil {
		return nil, err
	}
	api := API{
		client: customerservice.NewCustomerserviceClient(wiseEngageConn),
	}
	r := gin.New()
	r.Group("/callback").POST("/openim", api.OpenIMCallback) // Callback
	return nil, nil
}
