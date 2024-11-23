package worker

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/cralack/ChaosMetrics/server/internal/global"
	"github.com/cralack/ChaosMetrics/server/internal/service/pumper"
	"github.com/cralack/ChaosMetrics/server/utils"
	"github.com/cralack/ChaosMetrics/server/utils/register"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var workerID string
var podIP string
var HTTPListenAddress string
var GRPCListenAddress string
var PProfListenAddress string
var region string
var token string

var Cmd = &cobra.Command{
	Use:   "worker",
	Short: "start a pumper worker service",
	Args:  cobra.NoArgs,
	PreRun: func(cmd *cobra.Command, args []string) {
	},
	Run: func(cmd *cobra.Command, args []string) {
		Run(cmd.Context())
	},
}

func init() {
	Cmd.Flags().StringVar(&workerID, "id", "", "set worker id")
	Cmd.Flags().StringVar(&podIP, "podip", "192.168.123.197", "set pod IP")
	Cmd.Flags().StringVar(&HTTPListenAddress, "http", ":8082", "set HTTP listen address")
	Cmd.Flags().StringVar(&GRPCListenAddress, "grpc", ":9092", "set GRPC listen address")
	Cmd.Flags().StringVar(&PProfListenAddress, "pprof", ":9982", "set GRPC listen address")
	Cmd.Flags().StringVar(&region, "region", "AMERICAS", "set worker region")
	Cmd.Flags().StringVar(&token, "token", "", "set worker token")
}

func Run(ctx context.Context) {
	// load conf
	conf := global.ChaConf.System
	logger := global.ChaLogger
	conf.Name = global.WorkerServiceName
	conf.ID = workerID
	conf.GRPCListenAddress = GRPCListenAddress
	conf.HTTPListenAddress = HTTPListenAddress
	pumperID := conf.Name + "-" + conf.ID

	area := utils.ConvertRegionStrToArea(region)
	if workerID == "" {
		if podIP != "" {
			workerID = strconv.Itoa(int(utils.GetIDbyIP(podIP)))
		} else {
			workerID = fmt.Sprintf("%4d", time.Now().UnixNano())
		}
	}
	// init pumper core
	core, err := pumper.NewPumper(
		pumperID,
		pumper.WithAreaLoc(area),
		pumper.WithRegistryURL(conf.RegistryAddress),
		pumper.WithToken(token),
		pumper.WithContext(ctx),
	)
	if err != nil {
		logger.Panic("init worker failed", zap.Error(err))
		return
	}
	// start core
	core.StartEngine()
	logger.Info("starting worker engine...")

	go register.RunHTTPServer(logger, conf)
	register.RunGRPCServer(logger, conf)
}
