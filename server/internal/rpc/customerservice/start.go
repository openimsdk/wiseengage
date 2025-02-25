package customerservice

import (
	"context"

	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/wiseengage/v1/pkg/common/config"
	"github.com/openimsdk/wiseengage/v1/pkg/common/storage/controller"
	"github.com/openimsdk/wiseengage/v1/pkg/protocol/customerservice"
	"github.com/openimsdk/wiseengage/v1/pkg/rpcli"
	"google.golang.org/grpc"
)

type Config struct {
	Config config.Customer
}

func Start(ctx context.Context, config *Config, client discovery.SvcDiscoveryRegistry, server *grpc.Server) error {

	return nil
}

type customerService struct {
	customerservice.UnimplementedCustomerserviceServer
	db          controller.CustomerDatabase
	userCli     *rpcli.UserClient
	groupCli    *rpcli.GroupClient
	msgCli      *rpcli.MsgClient
	robotUserID string
	defaultRole string
	config      config.Customer
}
