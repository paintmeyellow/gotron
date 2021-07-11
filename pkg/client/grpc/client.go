package grpc

import (
	"github.com/paintmeyellow/tron-demo/api"
	"google.golang.org/grpc"
)

type Client struct {
	Address         string
	Conn            *grpc.ClientConn
	WalletClient    api.WalletClient
}

func NewClient(address string) *Client {
	var client Client
	client.Address = address
	return &client
}

func (c *Client) Start() error {
	var err error
	c.Conn, err = grpc.Dial(c.Address, grpc.WithInsecure())
	if err != nil {
		return err
	}
	c.WalletClient = api.NewWalletClient(c.Conn)
	return nil
}
