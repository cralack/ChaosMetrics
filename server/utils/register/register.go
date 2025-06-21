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
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go-micro.dev/v5"
	"go-micro.dev/v5/client"
	grpccli "go-micro.dev/v5/client/grpc"
	"go-micro.dev/v5/registry"
	"go-micro.dev/v5/registry/etcd"
	"go-micro.dev/v5/server"
	gs "go-micro.dev/v5/server/grpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Greeter struct{}

func (g *Greeter) Hello(ctx context.Context, in *emptypb.Empty, out *emptypb.Empty) error {
	_, cancel := context.WithCancel(ctx)
	defer cancel()
	if in == out {
		global.ChaLogger.Debug("in == out")
	}

	return nil
}

// RunGRPCServer and registry worker to etcd
func RunGRPCServer(logger *zap.Logger, cfg *config.MicroServ, opts ...interface{}) {
	// init grpc server
	reg := etcd.NewEtcdRegistry(registry.Addrs(cfg.RegistryAddress))

	service := micro.NewService(
		micro.Server(gs.NewServer(
			server.Name(cfg.Name),
			server.Id(cfg.ID),
			server.Address(cfg.GRPCListenAddress),
			server.Registry(reg),
			server.RegisterTTL(cfg.RegisterTTL*time.Second),
			server.RegisterInterval(cfg.RegisterInterval*time.Second),
		)),
		micro.Client(grpccli.NewClient()),
		// micro.WrapHandler(logWrapper(logger)),
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

func RunHTTPServer(logger *zap.Logger, cfg *config.MicroServ, opts ...interface{}) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mux := runtime.NewServeMux()
	dialOptions := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	_ = mux.HandlePath("GET", "/health", func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		w.WriteHeader(http.StatusOK)
		n, er := w.Write([]byte("ok"))
		if er != nil {
			logger.Error(er.Error())
		} else {
			logger.Debug(fmt.Sprintf("health response: %d", n))
		}
	})

	var err error
	if len(opts) > 0 {
		err = publisher.RegisterPublisherGwFromEndpoint(ctx, mux, cfg.GRPCListenAddress, dialOptions)
	} else {
		err = pb.RegisterGreeterGwFromEndpoint(ctx, mux, cfg.GRPCListenAddress, dialOptions)
	}

	if err != nil {
		logger.Fatal(fmt.Sprintf("register %s grpc http proxy failed", cfg.Name))
	}

	logger.Info(fmt.Sprintf("grpc's http proxy listening on %v", cfg.HTTPListenAddress))
	if err = http.ListenAndServe(cfg.HTTPListenAddress, mux); err != nil {
		logger.Fatal("HTTPListenAndServe failed")
	}
}

// func logWrapper(log *zap.Logger) server.HandlerWrapper {
// 	return func(hf server.HandlerFunc) server.HandlerFunc {
// 		return func(ctx context.Context, req server.Request, rsp interface{}) error {
// 			log.Info("receive request",
// 				zap.String("method", req.Method()),
// 				zap.String("Service", req.Service()),
// 				zap.Reflect("request param:", req.Body()),
// 			)
// 			err := hf(ctx, req, rsp)
// 			return err
// 		}
// 	}
// }
