package stripe

import (
	"github.com/stripe/stripe-go/v76/client"
)

type Client struct {
	client         *client.API
	endpointSecret string //whsec_GzeLxyXlwQI4mkYBhyWJWt5B5HFBfeOu
}

// NewClient returns a new stripe client
// stripeKey is the secret key for the stripe account
func (c *Client) NewClient(stripeKey string) *Client {
	return &Client{
		client: client.New(stripeKey, nil),
	}
}