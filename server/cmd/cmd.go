package cmd

import (
	"os"

	"github.com/cralack/ChaosMetrics/server/cmd/master"
	"github.com/cralack/ChaosMetrics/server/cmd/worker"
	"github.com/cralack/ChaosMetrics/server/global"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "run version service.",
	Long:  "run version service.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		global.GVA_LOG.Info("3.14.15")
	},
}

var rootCmd = &cobra.Command{
	Use:     "chao",
	Aliases: []string{"cm"},
	Short:   "chaosmetrics提供的命令行工具",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			global.GVA_LOG.Error("run root command help failed",
				zap.Error(err))
		}
	},
}

func AddCommands(root *cobra.Command) {
	root.AddCommand(
		master.Cmd,
		worker.Cmd,
		versionCmd,
	)
}

func RunCommand() error {
	AddCommands(rootCmd)
	{ // debug master part
		cmd, _, err := rootCmd.Find(os.Args[1:])
		if err != nil || cmd.Args == nil || global.GVA_ENV == global.TEST_ENV {
			arg := "master"
			extraArg1 := "--id=4"
			extraArg2 := "--http=:8084"
			extraArg3 := "--grpc=:9094"
			args := append([]string{arg, extraArg1, extraArg2, extraArg3}, os.Args[1:]...)
			rootCmd.SetArgs(args)
		}
	}
	// {
	// 	// debug worker part
	// 	cmd, _, err := rootCmd.Find(os.Args[1:])
	// 	if err != nil || cmd.Args == nil || global.GVA_ENV == global.TEST_ENV {
	// 		arg := "worker"
	// 		args := append([]string{arg}, os.Args[1:]...)
	// 		rootCmd.SetArgs(args)
	// 	}
	// }
	return rootCmd.Execute()
}
