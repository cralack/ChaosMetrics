package master

import (
	"github.com/cralack/ChaosMetrics/server/internal/global"
	"github.com/cralack/ChaosMetrics/server/service/master"
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
	conf := global.GvaConf.ServerConf
	logger := global.GvaLog
	conf.ID = masterId
	conf.Name = global.MasterServiceName
	conf.GRPCListenAddress = GRPCListenAddress
	conf.HTTPListenAddress = HTTPListenAddress

	m, err := master.New(
		masterId,
		master.WithLogger(logger.Named(global.MasterServiceName)),
		master.WithregistryURL(conf.RegistryAddress),
		master.WithGRPCAddress(conf.GRPCListenAddress),
		master.WithRegistry(etcd.NewRegistry(registry.Addrs(conf.RegistryAddress))),
	)
	if err != nil {
		logger.Error("start a master service failed", zap.Error(err))
	}

	m.Run()

	go register.RunHTTPServer(logger, conf)
	register.RunGRPCServer(logger, conf)
}
