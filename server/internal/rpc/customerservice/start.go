package customerservice

import (
	"context"

	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/wiseengage/v1/pkg/common/cmd"
	"github.com/openimsdk/wiseengage/v1/pkg/common/config"
	"github.com/openimsdk/wiseengage/v1/pkg/common/dbbuild"
	"github.com/openimsdk/wiseengage/v1/pkg/common/storage/controller"
	"github.com/openimsdk/wiseengage/v1/pkg/common/storage/database/mgo"
	"github.com/openimsdk/wiseengage/v1/pkg/imapi"
	"github.com/openimsdk/wiseengage/v1/pkg/protocol/customerservice"
	"google.golang.org/grpc"
)

type Config struct {
	Discovery     config.Discovery
	Config        config.Customer
	RedisConfig   config.Redis
	MongodbConfig config.Mongo
}

func Start(ctx context.Context, config *Config, client discovery.SvcDiscoveryRegistry) (cmd.BlockTasks, error) {
	dbb := dbbuild.NewBuilder(&config.MongodbConfig, &config.RedisConfig)
	mgocli, err := dbb.Mongo(ctx)
	if err != nil {
		return nil, err
	}
	customer, err := mgo.NewCustomer(mgocli.GetDB())
	if err != nil {
		return nil, err
	}
	srv := &customerService{
		db: controller.NewCustomerDatabase(customer, mgocli.GetTx()),
	}
	return func() error {
		return cmd.RPCServiceRegistrar(ctx, config.Discovery.RpcService.Customer, func(registrar grpc.ServiceRegistrar) {
			customerservice.RegisterCustomerserviceServer(registrar, srv)
		})
	}, nil
}

type customerService struct {
	customerservice.UnimplementedCustomerserviceServer
	db controller.CustomerDatabase
	//userCli  *rpcli.UserClient
	//groupCli *rpcli.GroupClient
	//msgCli      *rpcli.MsgClient
	//robotUserID string
	defaultRole string
	config      config.Customer
	imApi       imapi.CallerInterface
}
