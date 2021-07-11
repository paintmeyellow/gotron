package rest

import (
	"github.com/paintmeyellow/gotron/pkg/trongrid"
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
