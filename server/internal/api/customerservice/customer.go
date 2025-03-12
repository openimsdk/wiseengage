package customerapi

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/wiseengage/v1/pkg/protocol/customerservice"
)

func newGinRouter(ctx context.Context, client discovery.Conn, cfg *Config) (*gin.Engine, error) {
	customerConn, err := client.GetConn(ctx, cfg.Discovery.RpcService.Customer)
	if err != nil {
		return nil, err
	}
	api := API{
		client: customerservice.NewCustomerserviceClient(customerConn),
	}
	r := gin.New()
	r.Group("/callback").POST("/openim", api.OpenIMCallback) // Callback
	return nil, nil
}
