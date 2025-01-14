package stripe

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lgrisa/lib/utils/logutil"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/webhook"
)

func (c *Client) WebHookVerifySign(context *gin.Context) (*stripe.Event, error) {
	body, err := context.GetRawData()
	stripeSignature := context.GetHeader("Stripe-Signature")

	//验证签名
	event, err := webhook.ConstructEvent(body, stripeSignature, c.endpointSecret)
	if err != nil {
		return nil, err
	}

	logutil.LogTraceF("stripe webhook event.Type: %v", event.Type)

	if event.Data.Raw == nil {
		return nil, fmt.Errorf("event.Data.Raw is nil")
	}

	return &event, nil
}
