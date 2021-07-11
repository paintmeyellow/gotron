package grpc

import (
	"bytes"
	"context"
	"errors"
	"github.com/paintmeyellow/gotron/pkg/common/base58"
	"github.com/paintmeyellow/gotron/pkg/proto/core"
)

var ErrAccountNotFound = errors.New("account not found")

func (c *Client) GetAccount(ctx context.Context, addr string) (*core.Account, error) {
	var account core.Account
	var err error

	account.Address, err = base58.DecodeCheck(addr)
	if err != nil {
		return nil, err
	}

	acc, err := c.WalletClient.GetAccount(ctx, &account)
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(acc.Address, account.Address) {
		return nil, ErrAccountNotFound
	}
	return acc, nil
}