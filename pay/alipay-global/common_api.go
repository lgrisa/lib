package alipayGlobal

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lgrisa/lib/utils"
	"github.com/lgrisa/lib/utils/logutil"
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

	respHeader, httpCode, respBody, err := utils.Request(ctx, "POST", urlStr, header, bytes.NewBuffer([]byte(body)))

	if err != nil {
		return nil, err
	}

	if a.isDebugMode() {
		logutil.LogDebugF("Param：%s", body)

		logutil.LogDebugF("URL：%s", urlStr)

		logutil.LogDebugF("Resp Header：%v", respHeader)

		logutil.LogDebugF("Resp Body：%s", string(respBody))
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
