package cmd

import (
	"context"
	"os"

	"github.com/cralack/ChaosMetrics/server/cmd/master"
	"github.com/cralack/ChaosMetrics/server/cmd/worker"
	"github.com/cralack/ChaosMetrics/server/internal/global"
	"github.com/cralack/ChaosMetrics/server/internal/service/pumper"
	"github.com/cralack/ChaosMetrics/server/internal/service/router"
	"github.com/cralack/ChaosMetrics/server/internal/service/updater"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "get current env",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		switch global.ChaEnv {
		case global.TestEnv:
			global.ChaLogger.Info("enviroment:TEST")
		case global.ProductEnv:
			global.ChaLogger.Info("enviroment:PRODUCT")
		case global.DevEnv:
			global.ChaLogger.Info("enviroment:DEV")
		}
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "get current version",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		global.ChaLogger.Info("chaos metris current version: v0.9")
	},
}

var rootCmd = &cobra.Command{
	Use:     "chao",
	Aliases: []string{"cm"},
	Short:   "chaosmetrics提供的命令行工具",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			global.ChaLogger.Error("run root command help failed",
				zap.Error(err))
		}
	},
}

func AddCommands(root *cobra.Command) {
	root.AddCommand(
		master.Cmd,
		worker.Cmd,
		router.Cmd,
		pumper.Cmd,
		updater.Cmd,
		versionCmd,
		envCmd,
	)
}

func RunCmd(ctx context.Context) error {
	AddCommands(rootCmd)
	rootCmd.SetContext(ctx)

	{ // debug master part
		cmd, _, err := rootCmd.Find(os.Args[1:])
		if err != nil || cmd.Args == nil || global.ChaEnv == global.TestEnv {
			arg := "master"
			extraArg1 := "--id=2"
			extraArg2 := "--http=:8082"
			extraArg3 := "--grpc=:9092"
			args := append([]string{arg, extraArg1, extraArg2, extraArg3}, os.Args[1:]...)
			rootCmd.SetArgs(args)
		}
	}

	// { // debug worker part
	// 	cmd, _, err := rootCmd.Find(os.Args[1:])
	// 	if err != nil || cmd.Args == nil || global.ChaEnv == global.TestEnv {
	// 		arg := "worker"
	// 		// extraArg1 := "--id=1"
	// 		args := append([]string{arg}, os.Args[1:]...)
	// 		rootCmd.SetArgs(args)
	// 	}
	// }
	return rootCmd.Execute()
}
