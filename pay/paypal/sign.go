package paypal

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/lgrisa/library/utils"
	"github.com/pkg/errors"
	"net/http"
)

// VerifySign https://developer.paypal.com/docs/api/webhooks/v1/#verify-webhook-signature_post
func (c *Client) verifySign(context *gin.Context, webHookBody *[]byte) error {

	webHookeId := c.webhookId

	if webHookeId == "" {
		return errors.Errorf("VerifyPaypalSign MainConf.PayPalWebhookID == ''")
	}

	myVerifyReq := &verifyReq{
		AuthAlgo:         context.GetHeader("PAYPAL-AUTH-ALGO"),
		CertUrl:          context.GetHeader("PAYPAL-CERT-URL"),
		TransmissionId:   context.GetHeader("PAYPAL-TRANSMISSION-ID"),
		TransmissionSig:  context.GetHeader("PAYPAL-TRANSMISSION-SIG"),
		TransmissionTime: context.GetHeader("PAYPAL-TRANSMISSION-TIME"),
		WebhookId:        webHookeId,
		WebhookEvent:     json.RawMessage(*webHookBody),
	}

	verifyReqJson, err := json.Marshal(myVerifyReq)
	if err != nil {
		return errors.Errorf("VerifyPaypalSign json.Marshal(%v) error: %v", myVerifyReq, err)
	}

	urlProfile := httpProfile

	if !c.IsProd {
		urlProfile = httpSandboxProfile
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + c.accessToken,
	}

	_, respBody, err, httpCode := utils.HttpRequest(context, "POST", urlProfile+verifyMethod, headers, bytes.NewBuffer(verifyReqJson))

	if err != nil {
		return errors.Errorf("VerifyPaypalSign HttpRequest(%s,%s,%s) error: %v", "POST", urlProfile+verifyMethod, verifyReqJson, err)
	}

	if httpCode != 200 {
		return errors.Errorf("VerifyPaypalSign HttpRequest(%s,%s,%s) httpCode: %d", "POST", urlProfile+verifyMethod, verifyReqJson, httpCode)
	}

	myVerifyResp := &verifyResp{}
	err = json.Unmarshal(respBody, myVerifyResp)
	if err != nil {
		return errors.Errorf("VerifyPaypalSign json.Unmarshal(%v) error: %v", respBody, err)
	}

	if myVerifyResp.VerificationStatus != "SUCCESS" {
		return errors.Errorf("VerifyPaypalSign verifyResp.VerificationStatus != SUCCESS, verifyResp:%v", myVerifyResp)
	}

	return nil
}

func (c *Client) WebHookVerifySign(context *gin.Context) (*WebhookNotifyResponse, error) {
	rawData, err := context.GetRawData()
	if err != nil {
		return nil, err
	}

	if err = c.verifySign(context, &rawData); err != nil {
		return nil, err
	}

	resp := &WebhookNotifyResponse{}

	err = json.Unmarshal(rawData, resp)
	if err != nil {
		return nil, errors.Errorf("json.Unmarshal Err:%v", err)
	}

	return resp, nil
}

func (c *Client) WriteNotifyResp(context *gin.Context, code int) {
	if code == http.StatusOK {
		context.String(code, "%s", "success")
	} else {
		context.String(code, "%s", "fail")
	}
}
