package utils

import (
	"fmt"
	"github.com/awslabs/aws-lambda-go-api-proxy/core"
	"github.com/gin-gonic/gin"
	"github.com/lgrisa/lib/config"
	"github.com/lgrisa/lib/log"
	consts "github.com/lgrisa/lib/utils/const"
	"github.com/oschwald/geoip2-golang"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net"
	"os"
	"os/user"
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

func GetIpCurrencyCode(context *gin.Context) (currencyCode string) {
	var ok bool
	currencyMap := map[string]string{
		"CN": consts.CURRENCY_CNY,
		"US": consts.CURRENCY_USD,
		"JP": consts.CURRENCY_JPY,
		"KR": consts.CURRENCY_KRW,
		"PH": consts.CURRENCY_PHP,
		"TW": consts.CURRENCY_TWD,
		"MY": consts.CURRENCY_MYR,
		"HK": consts.CURRENCY_HKD,
		"SG": consts.CURRENCY_SGD,
		"TH": consts.CURRENCY_THB,
		"ID": consts.CURRENCY_IDR,
		"GB": consts.CURRENCY_GBP,
		"AU": consts.CURRENCY_AUD,
		"VN": consts.CURRENCY_VND,
		"CA": consts.CURRENCY_CAD,

		"CH": consts.CURRENCY_CHF, //瑞士
		"MO": consts.CURRENCY_MOP, //澳门
		"BR": consts.CURRENCY_BRL, //巴西
	}

	ipRegion := GetIpRegion(context)

	//欧元区
	euros := []string{"AU", "BE", "HR", "CY", "EE", "FI", "FR", "DE", "GR", "IE", "IT", "LV", "LT", "LU", "MT", "NL", "PT", "SK", "SI", "ES"}
	if InStringArray(ipRegion, euros, false) {
		return consts.CURRENCY_EUR
	}

	currencyCode, ok = currencyMap[ipRegion]
	if !ok {
		return consts.CURRENCY_USD
	}

	return
}

func GetIpRegion(context *gin.Context) (isoCode string) {
	clientIp := GetClientIp(context)
	//db, err := geoip2.Open(IpDbPath)

	//直接加载ip数据库，不使用缓存，因为geoip2库自己有缓存，而且我们的程序也不会频繁调用这个函数，所以不需要自己再加一层缓存
	db, err := geoip2.FromBytes(config.IpDbFs)

	if err != nil {
		log.LogErrorf("open ip db error: %v", err)
		return
	}
	defer db.Close()

	netIP := net.ParseIP(clientIp)
	record, err := db.City(netIP)
	if err != nil {
		log.LogErrorf("get ip region error: %v", err)
		return
	}

	isoCode = record.Country.IsoCode
	if isoCode == "" {
		isoCode = "CN"
	}
	return
}

func GetClientIp(context *gin.Context) (ipAddress string) {
	if context.GetHeader("iam-cf") == "true" {
		ipAddress = context.GetHeader("Zen-Client-Ip")
		if ipAddress != "" {
			return
		}

		xForwardedFor := context.GetHeader("X-Forwarded-For")
		arr := strings.Split(xForwardedFor, ",")
		if len(arr) > 0 {
			ipAddress = arr[0]
			return
		}
	} else {
		apiCtx, ok := core.GetAPIGatewayContextFromContext(context.Request.Context())
		if ok {
			ipAddress = apiCtx.Identity.SourceIP
		}
	}

	if ipAddress == "" {
		ipAddress = context.ClientIP()
	}

	return ipAddress
}

func GetUsername() string {
	username := os.Getenv("USER")

	if u, err := user.Current(); err == nil {
		if username != u.Username {
			username = fmt.Sprintf("%s(%s)", username, u.Username)
			if u.Username != u.Name {
				username = fmt.Sprintf("%s(%s)", username, u.Name)
			}
		} else {
			if username != u.Name {
				username = fmt.Sprintf("%s(%s)", username, u.Name)
			}
		}
	}

	return username
}

func GetHostname() string {
	hostname := os.Getenv("HOSTNAME")
	if hostname == "" {
		hostname = os.Getenv("HOST")
	}
	if hostname == "" {
		hostname, _ = os.Hostname()
	}
	if hostname == "" {
		hostname = "localhost"
	}
	return hostname
}
