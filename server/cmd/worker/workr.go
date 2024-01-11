package worker

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/cralack/ChaosMetrics/server/internal/config"
	"github.com/cralack/ChaosMetrics/server/internal/global"
	pb "github.com/cralack/ChaosMetrics/server/proto/greeter"
	"github.com/cralack/ChaosMetrics/server/service/pumper"
	"github.com/cralack/ChaosMetrics/server/utils"

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

var workerID string
var podIP string
var HTTPListenAddress string
var GRPCListenAddress string
var PProfListenAddress string
var cluster bool
var region string
var token string

var Cmd = &cobra.Command{
	Use:   "worker",
	Short: "start a pumper worker service",
	Args:  cobra.NoArgs,
	PreRun: func(cmd *cobra.Command, args []string) {
	},
	Run: func(cmd *cobra.Command, args []string) {
		Run()
	},
}

func init() {
	Cmd.Flags().StringVar(&workerID, "id", "", "set worker id")
	Cmd.Flags().StringVar(&podIP, "podip", "192.168.123.197", "set pod IP")
	Cmd.Flags().StringVar(&HTTPListenAddress, "http", ":8080", "set HTTP listen address")
	Cmd.Flags().StringVar(&GRPCListenAddress, "grpc", ":9090", "set GRPC listen address")
	Cmd.Flags().StringVar(&PProfListenAddress, "pprof", ":9981", "set GRPC listen address")
	Cmd.Flags().StringVar(&region, "region", "AMERICAS", "set worker region")
	Cmd.Flags().StringVar(&token, "token", "", "set worker token")
	Cmd.Flags().BoolVar(&cluster, "cluster", true, "run mode")
}

func Run() {
	// load conf
	conf := global.GVA_CONF.ServerConf
	logger := global.GVA_LOG

	area := utils.ConvertRegionToRegCode(region)
	if workerID == "" {
		if podIP != "" {
			workerID = strconv.Itoa(int(utils.GetIDbyIP(podIP)))
		} else {
			workerID = fmt.Sprintf("%4d", time.Now().UnixNano())
		}
	}

	// start pumper core
	exit := make(chan struct{})
	id := global.WorkerServiceName + "-" + workerID
	core, err := pumper.NewPumper(
		id,
		pumper.WithAreaLoc(area),
		pumper.WithRegistryURL(conf.RegistryAddress),
		pumper.WithToken(token),
	)
	if err != nil {
		logger.Panic("init worker failed", zap.Error(err))
		return
	}
	core.StartEngine(exit)
	logger.Info("starting worker engine...")

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

// RunGRPCServer and regestry worker to etcd
func RunGRPCServer(logger *zap.Logger, cfg *config.ServerConfig) {
	// init grpc server
	reg := etcdReg.NewRegistry(registry.Addrs(cfg.RegistryAddress))
	service := micro.NewService(
		micro.Server(gs.NewServer(server.Id(workerID))), // worker ID
		micro.Name(global.WorkerServiceName),            // worker name
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

	logger.Debug("worker grpc server starting")
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
				zap.String("method", req.Method()),
				zap.String("Service", req.Service()),
				zap.Reflect("request param:", req.Body()),
			)
			err := hf(ctx, req, rsp)
			return err
		}
	}
}
