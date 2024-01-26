package cmd

import (
	"github.com/cralack/ChaosMetrics/server/cmd/master"
	"github.com/cralack/ChaosMetrics/server/cmd/worker"
	"github.com/cralack/ChaosMetrics/server/internal/global"
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
		envCmd,
	)
}

func RunCmd() error {
	AddCommands(rootCmd)
	// { // debug master part
	// 	cmd, _, err := rootCmd.Find(os.Args[1:])
	// 	if err != nil || cmd.Args == nil || global.ChaEnv == global.TestEnv {
	// 		arg := "master"
	// 		extraArg1 := "--id=4"
	// 		extraArg2 := "--http=:8084"
	// 		extraArg3 := "--grpc=:9094"
	// 		args := append([]string{arg, extraArg1, extraArg2, extraArg3}, os.Args[1:]...)
	// 		rootCmd.SetArgs(args)
	// 	}
	// }
	// {
	// 	// debug worker part
	// 	cmd, _, err := rootCmd.Find(os.Args[1:])
	// 	if err != nil || cmd.Args == nil || global.ChaEnv == global.TestEnv {
	// 		arg := "worker"
	// 		extraArg1 := "--id=1"
	// 		args := append([]string{arg, extraArg1}, os.Args[1:]...)
	// 		rootCmd.SetArgs(args)
	// 	}
	// }
	return rootCmd.Execute()
}
