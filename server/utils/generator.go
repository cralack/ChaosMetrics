package utils

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math/rand"
	"net"

	"github.com/cralack/ChaosMetrics/server/internal/global"
	"go.uber.org/zap"
)

func GetIDbyIP(ip string) uint32 {
	global.ChaLogger.Debug(fmt.Sprintf("IP:[%s]", ip))
	var id uint32
	if err := binary.Read(bytes.NewBuffer(net.ParseIP(ip).To4()), binary.BigEndian, &id); err != nil {
		global.ChaLogger.Error("get id by ip failed",
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

func GenerateRandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var result = make([]byte, n)
	for i := 0; i < n; i++ {
		result[i] = letters[rand.Intn(len(result))]
	}
	return string(result)
}
