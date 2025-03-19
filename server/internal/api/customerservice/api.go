package customerapi

import (
	"context"

	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/wiseengage/v1/pkg/common/cmd"
	"github.com/openimsdk/wiseengage/v1/pkg/common/config"
	"github.com/openimsdk/wiseengage/v1/pkg/protocol/customerservice"
)

type Config struct {
	API       config.API
	Discovery config.Discovery
}

func Start(ctx context.Context, config *Config, client discovery.SvcDiscoveryRegistry) (cmd.BlockTasks, error) {
	router, err := newGinRouter(ctx, client, config)
	if err != nil {
		return nil, err
	}
	return func() error {
		return cmd.APIServer(ctx, router)
	}, nil
}

type API struct {
	client customerservice.CustomerserviceClient
}
