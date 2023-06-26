package test

import (
	"fmt"
	"os"
	"testing"

	_ "ChaosMetrics/server/init"
)

func Test_config(t *testing.T) {
	wd, _ := os.Getwd()
	fmt.Println(wd)
}
