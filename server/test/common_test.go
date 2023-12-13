package test

import (
	"testing"

	"github.com/cralack/ChaosMetrics/server/utils"
)

func Test_get_IP(t *testing.T) {
	if ip, err := utils.GetLocalIP(); err != nil {
		t.Log(err)
	} else {
		id := utils.GetIDbyIP(ip)
		t.Log(ip)
		t.Log(id)
	}

}
