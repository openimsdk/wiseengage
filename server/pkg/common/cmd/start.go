package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"reflect"
	"syscall"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mw"
	"github.com/openimsdk/tools/system/program"
	"github.com/openimsdk/wiseengage/v1/pkg/common/config"
	kdisc "github.com/openimsdk/wiseengage/v1/pkg/discovery"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type BlockTasks func() error

type RpcServer[S any] func(r grpc.ServiceRegistrar, srv S)

type StartFunc[C any] func(ctx context.Context, config *C, client discovery.SvcDiscoveryRegistry) (BlockTasks, error)

func Run[C any](fn StartFunc[C]) {
	cmd := cobra.Command{
		Use:           fmt.Sprintf("Start openIM application"),
		Long:          fmt.Sprintf(`Start %s `, program.GetProcessName()),
		SilenceUsage:  true,
		SilenceErrors: false,
		RunE: func(cmd *cobra.Command, args []string) error {
			configFolder, err := cmd.Flags().GetString("config_folder_path")
			if err != nil {
				return err
			}
			index, err := cmd.Flags().GetInt("index")
			if err != nil {
				return err
			}
			return run(fn, configFolder, index)
		},
	}
	if err := cmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
		return
	}
}

func run[C any](fn StartFunc[C], configFolder string, index int) error {
	conf, err := parseConfig[C](configFolder, index)
	if err != nil {
		return err
	}
	var discoveryConfig config.Discovery
	if err := readConfig(filepath.Join(configFolder, discoveryConfig.GetConfigFileName()), &discoveryConfig); err != nil {
		return err
	}
	client, err := kdisc.NewDiscoveryRegister(&discoveryConfig, nil)
	if err != nil {
		return err
	}
	defer client.Close()
	client.AddOption(
		mw.GrpcClient(), grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, "round_robin")),
	)

	rootCtx, exit := context.WithCancelCause(context.Background())
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
		select {
		case <-rootCtx.Done():
			return
		case val := <-sigs:
			log.ZDebug(rootCtx, "recv signal", "signal", val.String())
			exit(fmt.Errorf("signal %s", val.String()))
		}
	}()
	fnTimeout := func() (BlockTasks, error) {
		ctxValue := &startContextValue{
			RootCtx:      rootCtx,
			RootCancel:   exit,
			Config:       conf,
			Client:       client,
			Index:        index,
			ConfigFolder: configFolder,
		}
		ctx, cancel := context.WithTimeout(context.WithValue(rootCtx, startContextKey{}, ctxValue), time.Minute)
		defer cancel()
		return fn(ctx, conf, client)
	}
	block, err := fnTimeout()
	if err != nil {
		return err
	}
	if block == nil {
		return nil
	}
	waitBlock := make(chan error, 1)
	go func() {
		waitBlock <- block()
		close(waitBlock)
	}()
	select {
	case <-rootCtx.Done():
	case err := <-waitBlock:
		return err
	}
	ctxErr := context.Cause(rootCtx)
	timer := time.NewTimer(time.Second * 15)
	defer timer.Stop()
	select {
	case err := <-waitBlock:
		return err
	case <-timer.C:
		return fmt.Errorf("wait block timeout %w", ctxErr)
	}
}

func parseConfig[C any](configFolder string, index int) (*C, error) {
	var conf C
	vof := reflect.ValueOf(&conf)
	for vof.Kind() == reflect.Pointer {
		vof = vof.Elem()
	}
	if vof.Kind() != reflect.Struct {
		return nil, errors.New("config is not struct")
	}
	for i := 0; i < vof.NumField(); i++ {
		field := vof.Field(i)
		value := field.Interface()
		switch value.(type) {
		case config.Index:
			field.Set(reflect.ValueOf(config.Index(index)))
			continue
		case config.Path:
			field.Set(reflect.ValueOf(config.Path(configFolder)))
			continue
		}
		type ConfigFileName interface {
			GetConfigFileName() string
		}
		cfn, ok := value.(ConfigFileName)
		if !ok {
			return nil, fmt.Errorf("config %T not implement GetConfigFileName", value)
		}
		fieldValue := reflect.New(field.Type())
		if err := readConfig(filepath.Join(configFolder, cfn.GetConfigFileName()), fieldValue.Interface()); err != nil {
			return nil, err
		}
		field.Set(fieldValue.Elem())
	}
	return &conf, nil
}

func readConfig(name string, val any) error {
	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigFile(name)
	opt := func(conf *mapstructure.DecoderConfig) {
		conf.TagName = config.StructTagName
	}
	if err := v.Unmarshal(val, opt); err != nil {
		return err
	}
	return nil
}

type startContextKey struct{}

type startContextValue struct {
	RootCtx      context.Context
	RootCancel   context.CancelCauseFunc
	Config       any
	Client       discovery.SvcDiscoveryRegistry
	Index        int
	ConfigFolder string
}

func getStartContextValue(ctx context.Context) (*startContextValue, error) {
	value, ok := ctx.Value(startContextKey{}).(*startContextValue)
	if !ok {
		return nil, fmt.Errorf("invalid context")
	}
	return value, nil
}

func getSubConfig[C any](vof reflect.Value) *C {
	for vof.Kind() == reflect.Pointer {
		vof = vof.Elem()
	}
	if vof.Kind() != reflect.Struct {
		return nil
	}
	num := vof.NumField()
	for i := 0; i < num; i++ {
		field := vof.Field(i)
		for field.Kind() == reflect.Pointer {
			field = field.Elem()
		}
		if field.Kind() != reflect.Struct {
			continue
		}
		if !field.CanInterface() {
			continue
		}
		value := field.Interface()
		if rpc, ok := value.(C); ok {
			return &rpc
		}
		if rpc := getSubConfig[C](field); rpc != nil {
			return rpc
		}
	}
	return nil
}
