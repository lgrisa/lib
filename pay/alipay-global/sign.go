package alipayGlobal

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/packer/utils"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io"
	"net/url"
	"strings"
	"time"
)

// https://global.alipay.com/docs/ac/ams/api_fund#ML5ur
func (a *GlobalAliPayClient) getHeader(method, api, body string) (map[string]string, error) {
	timeStr := time.Now().Format("2006-01-02T15:04:05+08:00")

	sign, err := a.getSign(method, api, body, timeStr)
	if err != nil {
		return nil, err
	}

	signature := `algorithm=RSA256,keyVersion=1,signature=` + sign

	return map[string]string{
		"Signature":    signature,
		"client-id":    a.clientId,
		"request-time": timeStr,
		"Content-Type": "application/json; charset=UTF-8",
	}, nil
}

// https://global.alipay.com/docs/ac/ams/digital_signature#gNWs0
// generatedSignature=base64UrlEncode(sha256withrsa(<Content_To_Be_Signed>), <privateKey>))
func (a *GlobalAliPayClient) getSign(method, api, body, timeStr string) (sign string, err error) {

	signStr, err := utils.GetRsaSign(fmt.Sprintf("%s %s\n%s.%s.%s", method, api, a.clientId, timeStr, body), utils.RSA2, a.privateKey)

	if err != nil {
		return "", errors.Errorf("getSign GetRsaSign(%s,%s,%s,%s)：%v", method, api, body, timeStr, err)
	}

	return url.QueryEscape(signStr), nil
}

// https://global.alipay.com/docs/ac/ams/digital_signature#lt4nS
// IS_SIGNATURE_VALID=sha256withrsa_verify(base64UrlDecode(<target_signature>), <Content_To_Be_Validated>, <serverPublicKey>)
func (a *GlobalAliPayClient) VerifySign(ctx *gin.Context) error {
	data, err := ctx.GetRawData()
	if err != nil {
		return errors.Errorf("AliPay VerifySign ctx.GetRawData()：%v", err)
	}

	clientId := ctx.GetHeader("Client-Id")
	if clientId == "" {
		return errors.Errorf("AliPay VerifySign Client-Id is empty")
	}

	if clientId != a.clientId {
		return errors.Errorf("AliPay VerifySign Client-Id：%s", clientId)
	}

	reqTime := ctx.GetHeader("Request-Time")
	if reqTime == "" {
		return errors.Errorf("AliPay VerifySign Response-Time is empty")
	}

	signature := ctx.GetHeader("signature")
	if signature == "" {
		return errors.Errorf("AliPay VerifySign Signature is empty")
	}

	respSignature := strings.TrimPrefix(signature, "algorithm=RSA256,keyVersion=1,signature=")

	urlDecode, err := url.QueryUnescape(respSignature)
	if err != nil {
		return errors.Errorf("AliPay checkResp url.QueryUnescape(%s)：%v", respSignature, err)
	}

	checkSign := fmt.Sprintf("%s %s\n%s.%s.%s", ctx.Request.Method, ctx.Request.URL.Path, clientId, reqTime, string(data))

	logrus.Debugf("AliPay VerifySign checkSign:%s", checkSign)

	if err = utils.VerifySign(checkSign, urlDecode, utils.RSA2, a.aliPublicKey); err != nil {
		return errors.Errorf("AliPay VerifySign verifySign(%s,%s,%s)：%v", checkSign, urlDecode, utils.RSA2, err)
	}

	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(data))

	return nil
}
