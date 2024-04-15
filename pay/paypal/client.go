package paypal

import (
	"github.com/go-pay/gopay/paypal"
	"github.com/pkg/errors"
)

type Client struct {
	webhookId   string
	accessToken string
	ClientId    string
	Secret      string
	IsProd      bool

	GoPayClient *paypal.Client
}

func NewClient(webhookId, clientId, Secret string, isProd bool) (*Client, error) {

	goPayClient, err := paypal.NewClient(clientId, Secret, isProd)

	if err != nil {
		return nil, errors.Errorf("NewClient paypal.NewClient error: %v", err)
	}

	return &Client{
		webhookId:   webhookId,
		accessToken: goPayClient.AccessToken,
		ClientId:    clientId,
		Secret:      Secret,
		IsProd:      isProd,
		GoPayClient: goPayClient,
	}, nil
}
