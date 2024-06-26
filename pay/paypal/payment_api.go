package paypal

import (
	"context"
	"fmt"
	"github.com/go-pay/gopay/paypal"
	"github.com/lgrisa/lib/pay"
	"strings"
)

func (c *Client) EventCheckoutOrderApproved(ctx context.Context, notify *WebhookNotifyResponse) (*pay.CheckoutOrderApprovedResult, error) {
	if notify == nil {
		return nil, fmt.Errorf("notify is nil")
	}

	iapTransactionId := notify.Resource.Id

	if notify.EventType == WebhookEventCheckoutOrderApproved {

		ppRsp, err := c.GoPayClient.OrderCapture(ctx, iapTransactionId, nil)

		if err != nil {
			return nil, fmt.Errorf("PaymentAuthorizeCapture error: %v", err)
		}

		if ppRsp.Code != paypal.Success {
			if ppRsp.ErrorResponse != nil && ppRsp.ErrorResponse.Details != nil {
				if ppRsp.ErrorResponse.Details[0].Issue == "ORDER_ALREADY_CAPTURED" {
					return nil, pay.OrderAlreadyCapturedError
				} else {
					return nil, fmt.Errorf("PaymentAuthorizeCapture error: %v", ppRsp.ErrorResponse)
				}
			}
		}

		if ppRsp.Response.Status != "COMPLETED" {
			return nil, fmt.Errorf("PaymentAuthorizeCapture status: %v", ppRsp.Response.Status)
		}

		orderInfo := ppRsp.Response.PurchaseUnits[0]

		return &pay.CheckoutOrderApprovedResult{
			TransactionId:   orderInfo.ReferenceId,
			PlatformOrderId: iapTransactionId,
			RefundStr:       orderInfo.Payments.Captures[0].Id,
		}, nil

	} else {
		return nil, fmt.Errorf("event type is not CHECKOUT.ORDER.APPROVED")
	}
}

func (c *Client) EventCheckoutOrderComplete(ctx context.Context, notify *WebhookNotifyResponse) (*pay.CheckoutOrderApprovedResult, error) {
	if notify == nil {
		return nil, fmt.Errorf("notify is nil")
	}

	id := notify.Resource.SupplementaryData.RelatedIds["order_id"]

	if notify.EventType == WebhookEventCheckoutOrderComplete {

		ppRsp, err := c.GoPayClient.OrderDetail(ctx, id, nil)

		if err != nil {
			return nil, fmt.Errorf("OrderDetail error: %v", err)
		}

		if ppRsp.Code != paypal.Success {
			return nil, fmt.Errorf("OrderDetail ppRsp HttpStatusCode: %v,ErorrResponse: %v", ppRsp.Code, ppRsp.ErrorResponse)
		}

		if ppRsp.Response.Status != "COMPLETED" {
			return nil, fmt.Errorf("OrderDetail status: %v", ppRsp.Response.Status)
		}

		if len(ppRsp.Response.PurchaseUnits) > 1 {
			fmt.Println("PaypalNotify len(ppRsp.Response.PurchaseUnits) > 1")
		}

		return &pay.CheckoutOrderApprovedResult{
			TransactionId:   ppRsp.Response.PurchaseUnits[0].ReferenceId,
			PlatformOrderId: id,
			RefundStr:       ppRsp.Response.Id,
		}, nil

	} else {
		return nil, fmt.Errorf("event type is not CHECKOUT.ORDER.COMPLETE")
	}
}

func (c *Client) EventCheckoutOrderRefund(ctx context.Context, notify *WebhookNotifyResponse) (string, error) {
	if notify == nil {
		return "", fmt.Errorf("notify is nil")
	}

	id := notify.Resource.Id

	if notify.EventType == WebhookEventCheckoutOrderRefunded {

		ppRsp, err := c.GoPayClient.PaymentRefundDetail(ctx, id)

		if err != nil {
			return "", fmt.Errorf("PaymentRefundDetail error: %v", err)
		}

		if ppRsp.Code != paypal.Success {
			return "", fmt.Errorf("PaymentRefundDetail ppRsp HttpStatusCode: %v,ErorrResponse: %v", ppRsp.Code, ppRsp.ErrorResponse)
		}

		if ppRsp.Response.Status != "COMPLETED" {
			return "", fmt.Errorf("PaymentRefundDetail status: %v", ppRsp.Response.Status)
		}

		var refundStr string
		for _, v := range ppRsp.Response.Links {
			if v.Rel == "up" {
				list := strings.Split(v.Href, "/")

				refundStr = list[len(list)-1]
				break
			}
		}

		if refundStr == "" {
			return "", fmt.Errorf("PaymentRefundDetail refundStr is empty")
		}

		return refundStr, nil

	} else {
		return "", fmt.Errorf("event type is not CHECKOUT.ORDER.REFUND")
	}
}
