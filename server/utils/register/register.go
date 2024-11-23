package register

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/cralack/ChaosMetrics/server/internal/config"
	"github.com/cralack/ChaosMetrics/server/internal/global"
	"github.com/cralack/ChaosMetrics/server/internal/service/master"
	pb "github.com/cralack/ChaosMetrics/server/proto/greeter"
	"github.com/cralack/ChaosMetrics/server/proto/publisher"
	grpccli "github.com/go-micro/plugins/v4/client/grpc"

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
		global.ChaLogger.Debug("in == out")
	}

	return nil
}

// RunGRPCServer and regestry worker to etcd
func RunGRPCServer(logger *zap.Logger, cfg *config.SystemConf, opts ...interface{}) {
	// init grpc server
	reg := etcdReg.NewRegistry(registry.Addrs(cfg.RegistryAddress))
	service := micro.NewService(
		// start a grpc server
		micro.Server(gs.NewServer(server.Id(cfg.ID))),
		micro.Name(cfg.Name),
		micro.Address(cfg.GRPCListenAddress),
		micro.Registry(reg),
		micro.RegisterTTL(cfg.RegisterTTL*time.Second),
		micro.RegisterInterval(cfg.RegisterInterval*time.Second),
		micro.WrapHandler(logWrapper(logger)),
		// master forward require
		micro.Client(grpccli.NewClient()),
	)
	if err := service.Client().Init(client.RequestTimeout(cfg.ClientTimeOut * time.Second)); err != nil {
		logger.Error("micro client init error", zap.Error(err))
	}

	service.Init()

	var err error
	if len(opts) > 0 {
		m := opts[0].(*master.Master)
		cl := publisher.NewPublisherService(cfg.Name, service.Client())
		m.SetForwardCli(cl)

		err = publisher.RegisterPublisherHandler(service.Server(), m)
	} else {
		err = pb.RegisterGreeterHandler(service.Server(), new(Greeter))
	}
	if err != nil {
		logger.Fatal("register handler failed")
	}

	if err = service.Run(); err != nil {
		logger.Fatal(fmt.Sprintf("%s grpc server stop", cfg.Name))
	}
}

func RunHTTPServer(logger *zap.Logger, cfg *config.SystemConf, opts ...interface{}) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mux := runtime.NewServeMux()
	dialOptions := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	var err error
	if len(opts) > 0 {
		err = publisher.RegisterPublisherGwFromEndpoint(ctx, mux, cfg.GRPCListenAddress, dialOptions)
	} else {
		err = pb.RegisterGreeterGwFromEndpoint(ctx, mux, cfg.GRPCListenAddress, dialOptions)
	}

	if err != nil {
		logger.Fatal(fmt.Sprintf("register %s grpc http proxy failed", cfg.Name))
	}

	logger.Debug(fmt.Sprintf("grpc's http proxy listening on %v", cfg.HTTPListenAddress))
	if err = http.ListenAndServe(cfg.HTTPListenAddress, mux); err != nil {
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
