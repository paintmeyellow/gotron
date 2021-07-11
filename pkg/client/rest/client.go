package rest

import (
	"github.com/paintmeyellow/tron-demo/trongrid"
	"net/url"
)

type Client struct {
	WalletClient *trongrid.WalletClient
}

func NewClient(address url.URL) *Client {
	return &Client{
		WalletClient: trongrid.NewWalletClient(address),
	}
}
