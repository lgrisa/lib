package stripe

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lgrisa/library/utils"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/webhook"
)

const (
	PAY_COMPLETED       = "checkout.session.completed"
	PAY_ASYNC_SUCCEEDED = "checkout.session.async_payment_succeeded"
	REFUNDED            = "charge.refunded"
)

func (c *Client) WebHookVerify(context *gin.Context) (*stripe.Event, error) {
	body, err := context.GetRawData()
	stripeSignature := context.GetHeader("Stripe-Signature")

	//验证签名
	event, err := webhook.ConstructEvent(body, stripeSignature, c.endpointSecret)
	if err != nil {
		return nil, err
	}

	utils.LogTracef("stripe webhook event.Type: %v", event.Type)

	if event.Data.Raw == nil {
		return nil, fmt.Errorf("event.Data.Raw is nil")
	}

	return &event, nil
}
