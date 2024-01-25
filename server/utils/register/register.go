package register

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/cralack/ChaosMetrics/server/internal/config"
	"github.com/cralack/ChaosMetrics/server/internal/global"
	pb "github.com/cralack/ChaosMetrics/server/proto/greeter"
	etcdReg "github.com/go-micro/plugins/v4/registry/etcd"
	gs "github.com/go-micro/plugins/v4/server/grpc"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go-micro.dev/v4"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/server"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Greeter struct{}

func (g *Greeter) Hello(ctx context.Context, in *empty.Empty, out *empty.Empty) error {
	_, cancel := context.WithCancel(ctx)
	defer cancel()
	if in == out {
		global.GvaLog.Debug("in == out")
	}

	return nil
}

// RunGRPCServer and regestry worker to etcd
func RunGRPCServer(logger *zap.Logger, cfg *config.ServerConfig) {
	// init grpc server
	reg := etcdReg.NewRegistry(registry.Addrs(cfg.RegistryAddress))
	service := micro.NewService(
		micro.Server(gs.NewServer(server.Id(cfg.ID))),
		micro.Name(cfg.Name),
		micro.Address(cfg.GRPCListenAddress),
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
		logger.Fatal(fmt.Sprintf("%s grpc server stop", cfg.Name))
	}
}

func RunHTTPServer(logger *zap.Logger, cfg *config.ServerConfig) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := pb.RegisterGreeterGwFromEndpoint(ctx, mux, cfg.GRPCListenAddress, opts); err != nil {
		logger.Fatal(fmt.Sprintf("register %s grpc http proxy failed", cfg.Name))
	}

	logger.Debug(fmt.Sprintf("grpc's http proxy listening on %v", cfg.HTTPListenAddress))
	if err := http.ListenAndServe(cfg.HTTPListenAddress, mux); err != nil {
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
