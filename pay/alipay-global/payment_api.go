package alipayGlobal

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/lgrisa/library/utils"
	consts "github.com/lgrisa/library/utils/const"
	"github.com/pkg/errors"
	"net/http"
)

// FundTransRefund https://global.alipay.com/docs/ac/ams/refund_online
func (a *GlobalAliPayClient) FundTransRefund(ctx *gin.Context, bm utils.BodyMap) (*RefundResp, error) {
	err := bm.CheckEmptyError("transaction_id", "out_refund_no", "refund", "currency")
	if err != nil {
		return nil, err
	}

	req := &refundReq{
		PaymentId:       bm.Get("transaction_id"),
		RefundRequestId: bm.Get("out_refund_no"),
		RefundAmount: &refundAmount{
			Value:    int(bm.GetFloat64("refund") * 100),
			Currency: bm.GetString("currency"),
		},
	}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Errorf("json.Marshal(%v)：%v", req, err)
	}

	api := refundApi
	if !a.isProd {
		api = refundApiSandbox
	}

	resp, err := a.sendRequest(ctx, "POST", api, string(data))

	if err != nil {
		return nil, errors.Errorf("sendRequest(%s,%s,%s)：%s", "POST", api, data, err)
	}

	respJson := &RefundResp{}
	if err = json.Unmarshal(resp, respJson); err != nil {
		return nil, errors.Errorf("json.Unmarshal(%s)：%s", resp, err)
	}

	return respJson, nil
}

//https://global.alipay.com/docs/ac/ams/paymentrn_online

func (a *GlobalAliPayClient) WritePayNotifyResp(context *gin.Context, result string) error {
	resp := payCallBackResp{
		Result: struct {
			ResultCode    string `json:"resultCode"`
			ResultStatus  string `json:"resultStatus"`
			ResultMessage string `json:"resultMessage"`
		}{
			ResultCode:    "SUCCESS",
			ResultStatus:  "S",
			ResultMessage: "success",
		},
	}

	if result != "success" {
		resp.Result.ResultCode = "FAIL"
		resp.Result.ResultStatus = "F"
		resp.Result.ResultMessage = "Fail"
	}

	marshal, err := json.Marshal(resp)

	if err != nil {
		return errors.Errorf("json.Marshal(%v)：%v", resp, err)
	}

	context.Data(http.StatusOK, "application/json", marshal)

	return nil
}

//https://global.alipay.com/docs/ac/ams/notify_refund

func (a *GlobalAliPayClient) WriteAliRefundResp(context *gin.Context, result string) {
	resp := refundNotifyResp{
		Result: struct {
			ResultCode    string `json:"resultCode"`
			ResultStatus  string `json:"resultStatus"`
			ResultMessage string `json:"resultMessage"`
		}{
			ResultCode:    "SUCCESS",
			ResultStatus:  "S",
			ResultMessage: "Success",
		},
	}

	if result != "success" {
		resp.Result.ResultCode = "FAIL"
		resp.Result.ResultStatus = "F"
		resp.Result.ResultMessage = "Fail"
	}

	context.JSON(http.StatusOK, resp)
}

// https://global.alipay.com/docs/ac/ref/cc#ONkIe
func (a *GlobalAliPayClient) getAliPayAmount(amount float64, currencyCode string) int {
	if currencyCode != consts.CURRENCY_KRW {
		return int(amount * 100)
	} else {
		return int(amount)
	}
}

func (a *GlobalAliPayClient) payParamHandle(bm utils.BodyMap, RedirectUrl string, isMobile, isAndroid bool) (string, error) {
	payJson := &AliTradePagePay{
		ProductCode:      "CASHIER_PAYMENT",
		PaymentRequestId: bm.Get("out_trade_no"), //商户请求号
	}

	amountInt := a.getAliPayAmount(bm.GetFloat64("total_amount"), bm.GetString("currency"))

	payJson.Order.OrderAmount.Currency = bm.GetString("currency")
	payJson.Order.OrderAmount.Value = amountInt
	payJson.Order.ReferenceOrderId = bm.Get("out_trade_no")
	payJson.Order.OrderDescription = bm.Get("subject")

	payJson.PaymentAmount.Currency = bm.GetString("currency")
	payJson.PaymentAmount.Value = amountInt
	payJson.PaymentMethod.PaymentMethodType = "ALIPAY_CN"
	payJson.PaymentRedirectUrl = RedirectUrl

	payJson.PaymentNotifyUrl = a.webNotifyUrl
	payJson.SettlementStrategy.SettlementCurrency = bm.GetString("currency")
	payJson.Env.TerminalType = "WEB"

	if isMobile {
		payJson.Env.TerminalType = "WAP"

		if isAndroid {
			payJson.Env.OsType = "ANDROID"
		} else {
			payJson.Env.OsType = "IOS"
		}
	}

	marshal, err := utils.JSONMarshal(payJson)
	if err != nil {
		return "", errors.Errorf("payParamHandle json.Marshal(%v)：%v", payJson, err)
	}

	return string(marshal), nil
}

// TradePagePay https://global.alipay.com/docs/ac/ams/api_fund#WWH90
func (a *GlobalAliPayClient) TradePagePay(ctx *gin.Context, bm utils.BodyMap, returnUrl string, isMobile, isAndroid bool) (*PayResp, error) {
	sendParam, err := a.payParamHandle(bm, returnUrl, isMobile, isAndroid)
	if err != nil {
		return nil, errors.Errorf("payParamHandle(%v)：%v", bm, err)
	}

	api := payApi

	if !a.isProd {
		api = payApiSandbox
	}

	resp, err := a.sendRequest(ctx, "POST", api, sendParam)

	if err != nil {
		return nil, errors.Errorf("sendRequest(%s,%s,%s)：%s", "POST", api, sendParam, err)
	}

	payResp := &PayResp{}
	if err = json.Unmarshal(resp, payResp); err != nil {
		return nil, errors.Errorf("sendRequest json.Unmarshal(%s)：%s", resp, err)
	}

	if payResp.Result.ResultStatus == "S" {
	} else if payResp.Result.ResultStatus == "U" {
		if payResp.Result.ResultCode == "PAYMENT_IN_PROCESS" {

			payUrl := utils.NULL

			if payResp.SchemeUrl != "" {
				payUrl = payResp.SchemeUrl
			}

			if payResp.AppLinkUrl != "" && payUrl == utils.NULL {
				payUrl = payResp.AppLinkUrl
			}

			if payResp.NormalUrl != "" && payUrl == utils.NULL {
				payUrl = payResp.NormalUrl
			}

			if payUrl != utils.NULL {
			} else {
				return nil, errors.Errorf("sendRequest ResultCode：%s PaymentUrl is empty", payResp.Result.ResultCode)
			}

		} else {
			return nil, errors.Errorf("sendRequest ResultCode：%s ResultMessage：%s", payResp.Result.ResultCode, payResp.Result.ResultMessage)
		}
	}

	return payResp, nil
}

// TradeInquiryPayment https://global.alipay.com/docs/ac/ams/paymentri_online
func (a *GlobalAliPayClient) TradeInquiryPayment(ctx *gin.Context, paymentId string) (*InquiryPaymentResp, error) {

	if paymentId == "" {
		return nil, errors.Errorf("paymentId is empty")
	}

	req := InquiryPaymentReq{PaymentId: paymentId}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Errorf("json.Marshal(InquiryPaymentReq)：%v", err)
	}

	api := inquiryPaymentApi
	if !a.isProd {
		api = inquiryPaymentSandboxApi
	}

	resp, err := a.sendRequest(ctx, "POST", api, string(body))

	if err != nil {
		return nil, errors.Errorf("sendRequest(%s,%s,%s)：%s", "POST", api, body, err)
	}

	respJson := &InquiryPaymentResp{}
	if err = json.Unmarshal(resp, respJson); err != nil {
		return nil, errors.Errorf("sendRequest json.Unmarshal(%s)：%s", resp, err)
	}

	if respJson.Result.ResultStatus != "S" {
		return nil, errors.Errorf("sendRequest Result:%v", respJson.Result)
	}

	return respJson, nil
}
