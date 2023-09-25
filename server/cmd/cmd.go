package cmd

import (
	"os"

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

// matser entrance
var masterCmd = &cobra.Command{
	Use:   "master",
	Short: "start a master cluster",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
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
		masterCmd,
		worker.Cmd,
		versionCmd,
	)
}

func RunCommand() error {
	AddCommands(rootCmd)
	cmd, _, err := rootCmd.Find(os.Args[1:])
	if err != nil || cmd.Args == nil || global.GVA_ENV == global.TEST_ENV {
		// Not found
		args := append([]string{"worker"}, os.Args[1:]...)
		rootCmd.SetArgs(args)
	}
	return rootCmd.Execute()
}
