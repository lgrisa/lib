package alipayGlobal

import "time"

// refundReq 退款通知请求参数
type refundReq struct {
	PaymentId       string        `json:"paymentId"`
	RefundRequestId string        `json:"refundRequestId"`
	RefundAmount    *refundAmount `json:"refundAmount"`
}

type refundAmount struct {
	Value    int    `json:"value"`
	Currency string `json:"currency"`
}

// RefundResp 退款通知响应参数
type RefundResp struct {
	Result struct {
		ResultCode    string `json:"resultCode"`
		ResultStatus  string `json:"resultStatus"`
		ResultMessage string `json:"resultMessage"`
	} `json:"result"`
	RefundAmount struct {
		Value    string `json:"value"`
		Currency string `json:"currency"`
	} `json:"refundAmount"`
	RefundTime      time.Time `json:"refundTime"`
	PaymentId       string    `json:"paymentId"`
	RefundRequestId string    `json:"refundRequestId"`
	RefundId        string    `json:"refundId"`
}

type PayNotify struct {
	NotifyType string `json:"notifyType"`
	Result     struct {
		ResultCode    string `json:"resultCode"`
		ResultStatus  string `json:"resultStatus"`
		ResultMessage string `json:"resultMessage"`
	} `json:"result"`
	PaymentRequestId string `json:"paymentRequestId"`
	PaymentId        string `json:"paymentId"`
	PaymentAmount    struct {
		Value    string `json:"value"`
		Currency string `json:"currency"`
	} `json:"paymentAmount"`
	PaymentCreateTime time.Time `json:"paymentCreateTime"`
	PaymentTime       time.Time `json:"paymentTime"`
}

type payCallBackResp struct {
	Result struct {
		ResultCode    string `json:"resultCode"`
		ResultStatus  string `json:"resultStatus"`
		ResultMessage string `json:"resultMessage"`
	} `json:"result"`
}

type RefundNotify struct {
	NotifyType   string `json:"notifyType"`
	RefundAmount struct {
		Currency string `json:"currency"`
		Value    string `json:"value"`
	} `json:"refundAmount"`
	RefundId        string    `json:"refundId"`
	RefundRequestId string    `json:"refundRequestId"`
	RefundStatus    string    `json:"refundStatus"`
	RefundTime      time.Time `json:"refundTime"`
	Result          struct {
		ResultCode    string `json:"resultCode"`
		ResultMessage string `json:"resultMessage"`
		ResultStatus  string `json:"resultStatus"`
	} `json:"result"`
}

type refundNotifyResp struct {
	Result struct {
		ResultCode    string `json:"resultCode"`
		ResultStatus  string `json:"resultStatus"`
		ResultMessage string `json:"resultMessage"`
	} `json:"result"`
}

type AliTradePagePay struct {
	Order struct {
		OrderAmount struct {
			Currency string `json:"currency"`
			Value    int    `json:"value"`
		} `json:"orderAmount"`
		OrderDescription string `json:"orderDescription"`
		ReferenceOrderId string `json:"referenceOrderId"`
	} `json:"order"`
	Env struct {
		OsType       string `json:"osType"`
		TerminalType string `json:"terminalType"`
	} `json:"env"`
	PaymentAmount struct {
		Currency string `json:"currency"`
		Value    int    `json:"value"`
	} `json:"paymentAmount"`
	PaymentMethod struct {
		PaymentMethodType string `json:"paymentMethodType"`
	} `json:"paymentMethod"`
	SettlementStrategy struct {
		SettlementCurrency string `json:"settlementCurrency"`
	} `json:"settlementStrategy"`
	PaymentNotifyUrl   string `json:"paymentNotifyUrl"`
	PaymentRedirectUrl string `json:"paymentRedirectUrl"`
	PaymentRequestId   string `json:"paymentRequestId"`
	ProductCode        string `json:"productCode"`
}

type PayResp struct {
	AppIdentifier string `json:"appIdentifier"`
	AppLinkUrl    string `json:"applinkUrl"`
	NormalUrl     string `json:"normalUrl"`
	OrderCodeForm struct {
		CodeDetails []struct {
			CodeValue   string `json:"codeValue"`
			DisplayType string `json:"displayType"`
		} `json:"codeDetails"`
		ExpireTime time.Time `json:"expireTime"`
	} `json:"orderCodeForm"`
	PaymentAmount struct {
		Currency string `json:"currency"`
		Value    string `json:"value"`
	} `json:"paymentAmount"`
	PaymentCreateTime time.Time `json:"paymentCreateTime"`
	PaymentId         string    `json:"paymentId"` // 平台订单号
	PaymentRequestId  string    `json:"paymentRequestId"`
	SchemeUrl         string    `json:"schemeUrl"`
	Result            struct {
		ResultCode    string `json:"resultCode"`
		ResultMessage string `json:"resultMessage"`
		ResultStatus  string `json:"resultStatus"`
	} `json:"result"`
}

type InquiryPaymentReq struct {
	PaymentId string `json:"paymentId"`
}

type InquiryPaymentResp struct {
	Result struct {
		ResultCode    string `json:"resultCode"`
		ResultStatus  string `json:"resultStatus"`
		ResultMessage string `json:"resultMessage"`
	} `json:"result"`
	PaymentStatus    string `json:"paymentStatus"`
	PaymentRequestId string `json:"paymentRequestId"`
	PaymentId        string `json:"paymentId"`
	PaymentAmount    struct {
		Value    string `json:"value"`
		Currency string `json:"currency"`
	} `json:"paymentAmount"`
	PaymentCreateTime time.Time   `json:"paymentCreateTime"`
	PaymentTime       time.Time   `json:"paymentTime"`
	Transactions      interface{} `json:"transactions"`
}
