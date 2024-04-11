package paypal

import (
	"github.com/lgrisa/library/pay"
	"github.com/pkg/errors"
)

func (c *Client) EventCheckoutOrderApproved(resp *WebhookNotifyResponse) (*pay.CheckoutOrderApprovedResult, error) {
	if resp.EventType == WEBHOOK_EVENT_CHECKOUT_ORDER_APPROVED {
		//todo: Capture
	} else {
		return nil, errors.Errorf("event type is not CHECKOUT.ORDER.APPROVED")
	}

	return nil, nil
}
