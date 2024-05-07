package utils

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net"
	"strconv"
	"strings"
)

func ToU32Ip(addr [4]byte) uint32 {
	/* 127.0.0.1 = 0x100007f */
	var ip32 uint32
	ip32 |= uint32(addr[0])
	ip32 |= uint32(addr[1]) << 8
	ip32 |= uint32(addr[2]) << 16
	ip32 |= uint32(addr[3]) << 24

	return ip32
}

func ParseIPv4(ip string) ([]byte, error) {
	if ip == "" {
		return nil, errors.New("empty ip")
	}

	ip = strings.TrimSpace(ip)
	arr := strings.Split(ip, ".")
	if len(arr) != 4 {
		return nil, errors.New("invalid ip, " + ip)
	}

	var ipBytes [4]byte
	for i := 0; i < 4; i++ {
		v, err := strconv.Atoi(arr[i])
		if err != nil {
			return nil, errors.Wrap(err, "invalid ip, "+ip)
		}
		if v < 0 || v > 255 {
			return nil, errors.New("invalid ip, " + ip)
		}
		ipBytes[i] = byte(v)
	}

	return ipBytes[:], nil
}

// 解析本地地址
func ParseLocalAddr(localAddr string) []byte {
	var err error
	if localAddr == "" {
		localAddr, err = getLocalAddr()
		if err != nil {
			logrus.WithError(err).Panic("无法自动获得本机的内网地址, 请通过配置server.yaml中的local_addr字段手动设置")
		}
	}

	localAddrArray := make([]byte, 4)
	localAddrSplit := strings.Split(localAddr, ".")

	if len(localAddrSplit) != 4 {
		logrus.WithField("addr", localAddr).Panic("本地的机器地址local_addr必须是ip的形式. eg 192.168.1.10")
	}

	for i := 0; i < 4; i++ {
		num, err := strconv.Atoi(localAddrSplit[i])
		if err != nil || num < 0 || num >= 256 {
			logrus.WithField("addr", localAddr).Panic("本地的机器地址local_addr必须是个合法的ip. eg 192.168.1.10")
		}

		localAddrArray[i] = uint8(num)
	}
	return localAddrArray
}

func getLocalAddr() (string, error) {
	ips, err := getNetInterfaceIps()
	if err != nil {
		return "", err
	}

	return chooseBestLocalIp(ips), nil
}

func getNetInterfaceIps() ([]string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, errors.Wrap(err, "无法获得本地ip")
	}

	result := make([]string, 0, 2)

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			result = append(result, ipnet.IP.String())
		}
	}

	if len(result) == 0 {
		return nil, errors.Wrap(err, "没有从网卡上找到本地ip")
	}

	return result, nil
}

func chooseBestLocalIp(ips []string) string {

	for _, ip := range ips {
		if strings.HasPrefix(ip, "192.168") {
			return ip
		}
	}

	for _, ip := range ips {
		if strings.HasPrefix(ip, "172.16") {
			return ip
		}
	}

	for _, ip := range ips {
		if strings.HasPrefix(ip, "10.0") {
			return ip
		}
	}

	return ips[0]
}
