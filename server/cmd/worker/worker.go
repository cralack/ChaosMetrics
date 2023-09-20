package worker

import (
	"fmt"

	"github.com/cralack/ChaosMetrics/server/global"
	"github.com/cralack/ChaosMetrics/server/model/riotmodel"
	"github.com/cralack/ChaosMetrics/server/service/pumper"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var cluster bool
var workerID string
var podIP string
var HTTPListenAddress string
var GRPCListenAddress string
var PProfListenAddress string

var Cmd = &cobra.Command{
	Use:   "worker",
	Short: "start a pumper worker service",
	Args:  cobra.NoArgs,
	PreRun: func(cmd *cobra.Command, args []string) {

	},
	Run: func(cmd *cobra.Command, args []string) {
		global.GVA_LOG.Debug("workerid:",
			zap.String("flag", workerID))
		global.GVA_LOG.Debug("podip:",
			zap.String("flag", podIP))
		global.GVA_LOG.Debug("HTTPListenAddress:",
			zap.String("flag", HTTPListenAddress))
		global.GVA_LOG.Debug("GRPCListenAddress:",
			zap.String("flag", GRPCListenAddress))
		global.GVA_LOG.Debug("PProfListenAddress:",
			zap.String("flag", PProfListenAddress))
		global.GVA_LOG.Debug("cluster:",
			zap.Bool("flag", cluster))
		Run()
	},
}

func init() {
	Cmd.Flags().StringVar(
		&workerID, "id", "", "set worker id")
	Cmd.Flags().StringVar(
		&podIP, "podip", "", "set pod IP")
	Cmd.Flags().StringVar(
		&HTTPListenAddress, "http", ":8080", "set HTTP listen address")
	Cmd.Flags().StringVar(
		&GRPCListenAddress, "grpc", ":9090", "set GRPC listen address")
	Cmd.Flags().StringVar(
		&PProfListenAddress, "pprof", ":9981", "set GRPC listen address")
	Cmd.Flags().BoolVar(
		&cluster, "cluster", true, "run mode")
}

func Run() {
	// var (
	// 	err error
	// )
	// load conf
	conf := global.GVA_CONF.ServerConf
	conf.Name += ".worker"

	// logger := global.GVA_LOG
	// db := global.GVA_DB
	// rdb := global.GVA_RDB

	// start pumper core
	exit := make(chan struct{})
	core := pumper.NewPumper(
		pumper.WithLoc(riotmodel.NA1),
	)
	core.StartEngine(exit)

	fmt.Println("done")
}
