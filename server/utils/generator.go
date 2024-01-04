package utils

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net"

	"github.com/cralack/ChaosMetrics/server/internal/global"
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

// GetLocalIP 获取本机网卡IP
func GetLocalIP() (string, error) {
	var (
		addrs []net.Addr
		err   error
	)
	// get all network interface addr
	if addrs, err = net.InterfaceAddrs(); err != nil {
		return "", err
	}
	// get the first non-loopback network interface IP
	for _, addr := range addrs {
		if ipNet, isIpNet := addr.(*net.IPNet); isIpNet && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil
			}
		}
	}
	return "", errors.New("no local ip")
}
