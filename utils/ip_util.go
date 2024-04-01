package utils

import (
	"github.com/awslabs/aws-lambda-go-api-proxy/core"
	"github.com/gin-gonic/gin"
	"github.com/oschwald/geoip2-golang"
	"github.com/packer/utils/const"
	"github.com/sirupsen/logrus"
	"net"
	"strings"
)

func GetIpCurrencyCode(context *gin.Context, IpDbPath string) (currencyCode string) {
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

	ipRegion := GetIpRegion(context, IpDbPath)

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

func GetIpRegion(context *gin.Context, IpDbPath string) (isoCode string) {
	clientIp := GetClientIp(context)
	db, err := geoip2.Open(IpDbPath)
	if err != nil {
		logrus.Errorf("open ip db error: %v", err)
		return
	}
	defer db.Close()

	netIP := net.ParseIP(clientIp)
	record, err := db.City(netIP)
	if err != nil {
		logrus.Errorf("get ip region error: %v", err)
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
