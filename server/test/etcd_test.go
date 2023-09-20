package test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/cralack/ChaosMetrics/server/proto/greeter"
	"github.com/go-micro/plugins/v4/registry/etcd"
	gs "github.com/go-micro/plugins/v4/server/grpc"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Greeter struct{}

// var _ greeter.GreeterServer = &Greeter{}

func (g *Greeter) Hello(ctx context.Context, req *greeter.Request, rsp *greeter.Response) error {
	_, cancel := context.WithCancel(ctx)
	defer cancel()

	rsp.Greeting = "hello" + req.Name
	return nil
}

func Test_distrib(t *testing.T) {
	var err error
	// http proxy
	go HandleHTTP()

	reg := etcd.NewRegistry(
		registry.Addrs(":2379"),
	)
	// grpc server
	service := micro.NewService(
		micro.Server(gs.NewServer()),
		micro.Address(":9090"),
		micro.Registry(reg),
		micro.Name("pumper.worker"),
	)
	service.Init()
	if err = greeter.RegisterGreeterHandler(service.Server(), new(Greeter)); err != nil {
		fmt.Println(err)
	}

	if err = service.Run(); err != nil {
		fmt.Println(err)
	}
	fmt.Println("done")
}

func HandleHTTP() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := greeter.RegisterGreeterGwFromEndpoint(ctx, mux, "localhost:9090", opts)
	if err != nil {
		fmt.Println(err)
	}
	_ = http.ListenAndServe(":8080", mux)
}
