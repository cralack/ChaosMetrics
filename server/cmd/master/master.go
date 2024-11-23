package master

import (
	"context"

	"github.com/cralack/ChaosMetrics/server/internal/global"
	"github.com/cralack/ChaosMetrics/server/internal/service/master"
	"github.com/cralack/ChaosMetrics/server/utils/register"
	"github.com/go-micro/plugins/v4/registry/etcd"
	"github.com/spf13/cobra"
	"go-micro.dev/v4/registry"
	"go.uber.org/zap"
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
		Run(cmd.Context())
	},
}

func init() {
	Cmd.Flags().StringVar(&masterId, "id", "", "set master id")
	Cmd.Flags().StringVar(&podIP, "podip", "192.168.123.197", "set pod IP")
	Cmd.Flags().StringVar(&HTTPListenAddress, "http", ":8081", "set HTTP listen address")
	Cmd.Flags().StringVar(&GRPCListenAddress, "grpc", ":9091", "set GRPC listen address")
	Cmd.Flags().StringVar(&PProfListenAddress, "pprof", ":9981", "set GRPC listen address")
}

func Run(ctx context.Context) {
	conf := global.ChaConf.System
	logger := global.ChaLogger
	conf.ID = masterId
	conf.Name = global.MasterServiceName
	conf.GRPCListenAddress = GRPCListenAddress
	conf.HTTPListenAddress = HTTPListenAddress

	// init master
	m, err := master.New(
		conf.Name+"-"+conf.ID,
		master.WithLogger(logger.Named(global.MasterServiceName)),
		master.WithregistryURL(conf.RegistryAddress),
		master.WithGRPCAddress(conf.GRPCListenAddress),
		master.WithRegistry(etcd.NewRegistry(registry.Addrs(conf.RegistryAddress))),
		master.WithContext(ctx),
	)
	if err != nil {
		logger.Error("start a master service failed", zap.Error(err))
	}

	m.Run()

	go register.RunHTTPServer(logger, conf, m)
	register.RunGRPCServer(logger, conf, m)
}
