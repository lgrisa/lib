package alipayGlobal

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/packer/utils"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"strings"
)

func (a *GlobalAliPayClient) sendRequest(ctx *gin.Context, method, api, body string) ([]byte, error) {

	header, err := a.getHeader(method, api, body)
	if err != nil {
		return nil, err
	}

	urlStr := AliUrl + api

	if !a.isProd {
		urlStr = AliUrlSandbox + api
	}

	respHeader, respBody, err, httpCode := utils.HttpRequest(ctx, "POST", urlStr, header, bytes.NewBuffer([]byte(body)))

	if err != nil {
		return nil, err
	}

	if a.isDebugMode() {
		utils.LogDebugf("Param：%s", body)

		utils.LogDebugf("URL：%s", urlStr)

		utils.LogDebugf("Resp Header：%v", respHeader)

		utils.LogDebugf("Resp Body：%s", string(respBody))
	}

	if httpCode != 200 {
		return nil, fmt.Errorf("HTTP Request Error, StatusCode = %d", httpCode)
	}

	isOk, err := a.checkResp(&respHeader, method, api, &respBody)

	if err != nil {
		return nil, fmt.Errorf("checkResp err：%v", err)
	} else if !isOk {
		return nil, fmt.Errorf("checkResp is not ok")
	}

	return respBody, nil
}

// https://global.alipay.com/docs/ac/ams/api_fund#WWH90
func (a *GlobalAliPayClient) checkResp(respHeader *http.Header, method, api string, respBody *[]byte) (bool, error) {

	respSign := respHeader.Get("signature")
	if respSign == "" {
		return false, errors.Errorf("Signature is empty")
	}

	respSignature := strings.TrimPrefix(respSign, "algorithm=RSA256,keyVersion=1,signature=")

	clientId := respHeader.Get("client-id")
	if clientId != a.clientId {
		return false, errors.Errorf("clientId：%s", clientId)
	}

	timeStr := respHeader.Get("Response-Time")
	if timeStr == "" {
		return false, errors.Errorf("Response-Time is empty")
	}

	checkBody := fmt.Sprintf("%s %s\n%s.%s.%s", method, api, a.clientId, timeStr, string(*respBody))

	urlDecode, err := url.QueryUnescape(respSignature)
	if err != nil {
		return false, errors.Errorf("url.QueryUnescape(%s)：%v", respSignature, err)
	}

	if err = utils.VerifySign(checkBody, urlDecode, utils.RSA2, a.aliPublicKey); err != nil {
		return false, err
	}

	return true, nil
}