package stripe

import (
	"encoding/json"
	"fmt"
	"github.com/lgrisa/lib/pay"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/checkout/session"
)

// CreateOrder creates a new order https://docs.stripe.com/api/checkout/sessions/create
// priceId Background priceId
// outTradeNo The order number of the recharge
// successUrl The URL to redirect to after the transaction is successful
// transactionUrl The URL to redirect to after the transaction is successful
// platFromOrderId The order number of the recharge
func (c *Client) CreateOrder(priceId string, outTradeNo string, successUrl string) (string, string, error) {

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
		return "", "", err
	}

	return curSession.URL, curSession.ID, nil
}

func (c *Client) CheckoutSessionCompleted(event *stripe.Event) (*pay.CheckoutOrderApprovedResult, error) {

	if event.Type != PayCompleted && event.Type != PayAsyncSucceeded {
		return nil, fmt.Errorf("event.Type != checkout.session.completed")
	}

	checkoutSession := &stripe.CheckoutSession{}
	err := json.Unmarshal(event.Data.Raw, checkoutSession)

	if err != nil {
		return nil, err
	}

	if checkoutSession.PaymentStatus == stripe.CheckoutSessionPaymentStatusPaid {
		return &pay.CheckoutOrderApprovedResult{
			TransactionId:   checkoutSession.ClientReferenceID,
			PlatformOrderId: checkoutSession.ID,
			RefundStr:       checkoutSession.PaymentIntent.ID,
		}, nil
	} else {
		return nil, fmt.Errorf("checkoutSession.PaymentStatus == %v", checkoutSession.PaymentStatus)
	}
}

// CheckoutSessionRefunded https://docs.stripe.com/api/refunds/create
func (c *Client) CheckoutSessionRefunded(event *stripe.Event) (string, error) {
	if event.Type != REFUNDED {
		return "", fmt.Errorf("event.Type != charge.refunded")
	}

	refund := &stripe.Refund{}
	if err := json.Unmarshal(event.Data.Raw, refund); err != nil {
		return "", err
	}

	if refund.Status == stripe.RefundStatusSucceeded {
		return refund.PaymentIntent.ID, nil
	} else {
		return "", fmt.Errorf("refund.Status == %v", refund.Status)
	}
}
