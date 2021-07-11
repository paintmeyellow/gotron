package grpc

import (
	"context"
	"crypto/ecdsa"
	"github.com/paintmeyellow/gotron/pkg/common/crypto"
	"github.com/paintmeyellow/gotron/pkg/common/types"
	"github.com/paintmeyellow/gotron/pkg/proto/core"
	"github.com/pkg/errors"
)

var ErrInvalidTransactionResult = errors.New("invalid transaction result")

type SendableTransaction struct {
	OwnerKey  ecdsa.PrivateKey
	ToAddress *crypto.Address
	Amount    int64
}

func (c *Client) Transfer(ctx context.Context, sendTxn SendableTransaction) (*crypto.Hash, error) {
	contract := core.TransferContract{
		OwnerAddress: crypto.PubkeyToAddress(sendTxn.OwnerKey.PublicKey).Bytes(),
		ToAddress:    sendTxn.ToAddress.Bytes(),
		Amount:       sendTxn.Amount,
	}

	txn, err := c.WalletClient.CreateTransaction2(ctx, &contract)
	if err != nil {
		return nil, err
	}

	if txn == nil || txn.GetResult().GetCode() != 0 {
		return nil, ErrInvalidTransactionResult
	}

	if err = types.SignTxn(txn, &sendTxn.OwnerKey); err != nil {
		return nil, err
	}

	_, err = c.WalletClient.BroadcastTransaction(ctx, txn.GetTransaction())
	if err != nil {
		return nil, err
	}

	var hash crypto.Hash
	hash.SetBytes(txn.GetTxid())

	return &hash, nil
}
