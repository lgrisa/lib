package stripe

import (
	"encoding/json"
	"fmt"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
	"github.com/lgrisa/library/pay"
)

// CreateOrder creates a new order https://docs.stripe.com/api/checkout/sessions/create
// priceId Background priceId
// outTradeNo The order number of the recharge
// successUrl The URL to redirect to after the transaction is successful
// transactionUrl The URL to redirect to after the transaction is successful
// platFromOrderId The order number of the recharge
func (c *Client) CreateOrder(priceId string, outTradeNo string, successUrl string) (transactionUrl string, platFromOrderId string) {

	params := &stripe.CheckoutSessionParams{
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(priceId),
				Quantity: stripe.Int64(1),
			},
		},
		Mode:              stripe.String("payment"),
		SuccessURL:        stripe.String(successUrl),
		ClientReferenceID: stripe.String(outTradeNo),
	}

	if successUrl != "" {
		params.SuccessURL = stripe.String(successUrl)
	}

	curSession, err := session.New(params)
	if err != nil {
		return
	}

	transactionUrl = curSession.URL

	platFromOrderId = curSession.ID

	return
}

func (c *Client) CheckoutSessionCompleted(rawData []byte) (*pay.CheckoutOrderApprovedResult, error) {
	checkoutSession := &stripe.CheckoutSession{}
	err := json.Unmarshal(rawData, checkoutSession)
	if err != nil {
		return nil, err
	}

	if checkoutSession.PaymentStatus == stripe.CheckoutSessionPaymentStatusPaid {
		return &pay.CheckoutOrderApprovedResult{
			OutTradeNo:      checkoutSession.ClientReferenceID,
			PlatformOrderId: checkoutSession.ID,
			RefundStr:       checkoutSession.PaymentIntent.ID,
		}, nil
	} else {
		return nil, fmt.Errorf("checkoutSession.PaymentStatus == %v", checkoutSession.PaymentStatus)
	}
}
