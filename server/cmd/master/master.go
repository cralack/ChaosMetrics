package master

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/cralack/ChaosMetrics/server/internal/config"
	"github.com/cralack/ChaosMetrics/server/internal/global"
	pb "github.com/cralack/ChaosMetrics/server/proto/greeter"
	"github.com/cralack/ChaosMetrics/server/service/master"
	"github.com/go-micro/plugins/v4/registry/etcd"
	etcdReg "github.com/go-micro/plugins/v4/registry/etcd"
	gs "github.com/go-micro/plugins/v4/server/grpc"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/cobra"
	"go-micro.dev/v4"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/server"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var podIP string
var masterId string
var HTTPListenAddress string
var GRPCListenAddress string
var PProfListenAddress string

var Cmd = &cobra.Command{
	Use:   "master",
	Short: "start a master cluster",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		Run()
	},
}

func init() {
	Cmd.Flags().StringVar(&masterId, "id", "", "set master id")
	Cmd.Flags().StringVar(&podIP, "podip", "192.168.123.197", "set pod IP")
	Cmd.Flags().StringVar(&HTTPListenAddress, "http", ":8080", "set HTTP listen address")
	Cmd.Flags().StringVar(&GRPCListenAddress, "grpc", ":9090", "set GRPC listen address")
	Cmd.Flags().StringVar(&PProfListenAddress, "pprof", ":9981", "set GRPC listen address")
	Cmd.Flags().StringVar(&PProfListenAddress, "token", "", "set riot API token")
}

func Run() {
	conf := global.GVA_CONF.ServerConf
	logger := global.GVA_LOG
	reg := etcd.NewRegistry(registry.Addrs(conf.RegistryAddress))

	m, err := master.New(
		masterId,
		master.WithLogger(logger.Named("master")),
		master.WithregistryURL(conf.RegistryAddress),
		master.WithGRPCAddress(GRPCListenAddress),
		master.WithRegistry(reg),
	)
	if err != nil {
		logger.Error("start a master service failed", zap.Error(err))
	}

	m.Run()

	go RunHTTPServer(logger)
	RunGRPCServer(logger, conf)
}

type Greeter struct{}

func (g *Greeter) Hello(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
	_, cancel := context.WithCancel(ctx)
	defer cancel()

	rsp.Greeting = "Hello," + req.Name
	return nil
}

// RunGRPCServer and regestry master to etcd
func RunGRPCServer(logger *zap.Logger, cfg *config.ServerConfig) {
	// init grpc server
	reg := etcdReg.NewRegistry(registry.Addrs(cfg.RegistryAddress))
	service := micro.NewService(
		micro.Server(gs.NewServer(server.Id(masterId))),
		micro.Name(global.MasterServiceName),
		micro.Address(GRPCListenAddress),
		micro.Registry(reg),
		micro.RegisterTTL(cfg.RegisterTTL*time.Second),
		micro.RegisterInterval(cfg.RegisterInterval*time.Second),
		micro.WrapHandler(logWrapper(logger)),
	)
	if err := service.Client().Init(client.RequestTimeout(cfg.ClientTimeOut * time.Second)); err != nil {
		logger.Error("micro client init error", zap.Error(err))
	}

	service.Init()

	if err := pb.RegisterGreeterHandler(service.Server(), new(Greeter)); err != nil {
		logger.Fatal("register handler failed")
	}

	if err := service.Run(); err != nil {
		logger.Fatal("worker grpc server stop")
	}
}

func RunHTTPServer(logger *zap.Logger) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := pb.RegisterGreeterGwFromEndpoint(ctx, mux, GRPCListenAddress, opts); err != nil {
		logger.Fatal("register worker grpc http proxy failed")
	}

	logger.Debug(fmt.Sprintf("grpc's http proxy listening on %v", HTTPListenAddress))
	if err := http.ListenAndServe(HTTPListenAddress, mux); err != nil {
		logger.Fatal("HTTPListenAndServe failed")
	}
}

func logWrapper(log *zap.Logger) server.HandlerWrapper {
	return func(hf server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			log.Info("receive request",
				zap.String("time", time.Now().Format(time.RFC3339)),
				zap.String("method", req.Method()),
				zap.String("Service", req.Service()),
				zap.Reflect("request param:", req.Body()),
			)
			err := hf(ctx, req, rsp)
			return err
		}
	}
}
