package paypal

const (
	httpProfile        = "https://api-m.paypal.com"
	httpSandboxProfile = "https://api-m.sandbox.paypal.com"

	verifyMethod = "/v1/notifications/verify-webhook-signature"
	orderCapture = "/v2/checkout/orders/%s/capture"
)

const (
	WebhookEventCheckoutOrderApproved = "CHECKOUT.ORDER.APPROVED"
	WebhookEventCheckoutOrderComplete = "PAYMENT.CAPTURE.COMPLETED"
	WebhookEventCheckoutOrderRefunded = "PAYMENT.CAPTURE.REFUNDED"
)
