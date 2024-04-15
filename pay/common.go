package pay

import "github.com/pkg/errors"

type DebugSwitch int8

const (
	DebugOff = 0
	DebugOn  = 1
)

type CheckoutOrderApprovedResult struct {
	TransactionId   string `json:"out_trade_no"`      // 商户订单号
	PlatformOrderId string `json:"platform_order_id"` // 平台订单号
	RefundStr       string `json:"refund_str"`        // 退款字符串
}

var OrderAlreadyCapturedError = errors.New("order already captured")
