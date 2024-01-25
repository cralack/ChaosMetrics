package worker

import (
	"fmt"
	"strconv"
	"time"

	"github.com/cralack/ChaosMetrics/server/internal/global"
	"github.com/cralack/ChaosMetrics/server/service/pumper"
	"github.com/cralack/ChaosMetrics/server/service/register"
	"github.com/cralack/ChaosMetrics/server/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
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
	conf := global.GvaConf.ServerConf
	logger := global.GvaLog
	conf.Name = global.WorkerServiceName
	conf.ID = workerID
	conf.GRPCListenAddress = GRPCListenAddress
	conf.HTTPListenAddress = HTTPListenAddress

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
	core, err := pumper.NewPumper(
		conf.ID,
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

	go register.RunHTTPServer(logger, conf)
	register.RunGRPCServer(logger, conf)
}
