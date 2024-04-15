package stripe

//https://docs.stripe.com/api/events/types#event_types-checkout.session.completed

const (
	PayCompleted      = "checkout.session.completed"
	PayAsyncSucceeded = "checkout.session.async_payment_succeeded"
	REFUNDED          = "charge.refunded"
)
