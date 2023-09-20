package utils

import (
	"bytes"
	"encoding/binary"
	"net"

	"github.com/cralack/ChaosMetrics/server/global"
	"go.uber.org/zap"
)

func GetIDbyIP(ip string) uint32 {
	var id uint32
	if err := binary.Read(bytes.NewBuffer(net.ParseIP(ip).To4()), binary.BigEndian, &id); err != nil {
		global.GVA_LOG.Error("get id by ip failed",
			zap.String("utils", err.Error()))
	}
	return id
}
