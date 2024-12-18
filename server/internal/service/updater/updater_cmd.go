package updater

import (
	"time"

	"github.com/spf13/cobra"
)

var (
	force   bool
	endmark string
)
var Cmd = &cobra.Command{
	Use:   "update",
	Short: "start a local updater process",
	Run: func(cmd *cobra.Command, args []string) {
		Run()
	},
}

func init() {
	Cmd.Flags().BoolVar(&force, "force", false, "if force update")
	Cmd.Flags().StringVarP(&endmark, "endmark", "e", "14.1.1", "version endmark")
}

func Run() {
	u := NewRiotUpdater(
		WithLifeTime(time.Hour*24*30),
		WithForceUpdate(force),
		WithEndmark(endmark),
	)

	u.UpdateAll()
}
