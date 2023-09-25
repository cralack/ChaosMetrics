package main

import (
	"github.com/cralack/ChaosMetrics/server/cmd"
	_ "github.com/cralack/ChaosMetrics/server/init"
)

func main() {
	if err := cmd.RunCommand(); err != nil {
		panic(err)
	}
}
