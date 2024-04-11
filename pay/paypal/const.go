package paypal

const (
	verifyUrl        = "https://api-m.paypal.com/v1/notifications/verify-webhook-signature"
	verifySandboxUrl = "https://api-m.sandbox.paypal.com/v1/notifications/verify-webhook-signature"
)

const(
	WEBHOOK_EVENT_CHECKOUT_ORDER_APPROVED = "CHECKOUT.ORDER.APPROVED"
)