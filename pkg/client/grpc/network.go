package grpc

import (
	"bytes"
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/paintmeyellow/tron-demo/api"
	"github.com/paintmeyellow/tron-demo/common/crypto"
	"github.com/paintmeyellow/tron-demo/common/hexutil"
	"github.com/paintmeyellow/tron-demo/core"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"time"
)

var (
	ErrTransactionNotFound    = errors.New("transaction not found")
	ErrInvalidBroadcastResult = errors.New("invalid broadcast transaction result")
	ErrEmptyContractList = errors.New("contract list is empty")
)

type Transaction struct {
	Hash   *crypto.Hash
	From   *crypto.Address
	To     *crypto.Address
	NetFee int64
	Amount int64
}

//GetTransactionInfoByID returns transaction receipt by ID
func (c *Client) GetTransactionInfoByID(ctx context.Context, hash *crypto.Hash) (*core.TransactionInfo, error) {
	var txnID api.BytesMessage
	txnID.Value = hash.Bytes()
	txn, err := c.WalletClient.GetTransactionInfoById(ctx, &txnID)
	if err != nil {
		return nil, err
	}
	if bytes.Equal(txn.Id, txnID.Value) {
		return txn, nil
	}
	return nil, ErrTransactionNotFound
}

//GetTransactionByID returns transaction details by ID
func (c *Client) GetTransactionByID(ctx context.Context, hash *crypto.Hash) (*core.Transaction, error) {
	var txnID api.BytesMessage
	txnID.Value = hash.Bytes()
	txn, err := c.WalletClient.GetTransactionById(ctx, &txnID)
	if err != nil {
		return nil, err
	}
	if proto.Size(txn) > 0 {
		return txn, nil
	}
	return nil, ErrTransactionNotFound
}

//GetTransactionFullInfoByID returns transaction details by ID
func (c *Client) GetTransactionFullInfoByID(ctx context.Context, h *crypto.Hash, tries uint, delay time.Duration) (*Transaction, error) {
	var transaction Transaction
	errs, ctx := errgroup.WithContext(ctx)

	errs.Go(func() error {
		txn, err := c.GetTransactionByID(ctx, h)
		if err != nil {
			return err
		}
		if len(txn.GetRawData().GetContract()) == 0 {
			return ErrEmptyContractList
		}
		var contract core.TransferContract
		if err = txn.GetRawData().GetContract()[0].GetParameter().UnmarshalTo(&contract); err != nil {
			return err
		}
		transaction.From = crypto.HexToAddress(hexutil.Bytes2Hex(contract.OwnerAddress))
		transaction.To = crypto.HexToAddress(hexutil.Bytes2Hex(contract.ToAddress))
		transaction.Amount = contract.Amount
		return nil
	})

	errs.Go(func() error {
		for tries > 0 {
			<-time.Tick(delay)
			txn, err := c.GetTransactionInfoByID(ctx, h)
			if errors.Is(err, ErrTransactionNotFound) {
				tries--
				continue
			}
			if err != nil {
				return err
			}
			transaction.Hash = crypto.BytesToHash(txn.Id)
			transaction.NetFee = txn.Fee
			return nil
		}
		return ErrTransactionNotFound
	})

	return &transaction, errs.Wait()
}

func (c *Client) Broadcast(ctx context.Context, tx *core.Transaction) (*api.Return, error) {
	result, err := c.WalletClient.BroadcastTransaction(ctx, tx)
	if err != nil {
		return nil, err
	}
	if !result.GetResult() || result.GetCode() != api.Return_SUCCESS {
		return result, errors.Wrapf(ErrInvalidBroadcastResult,
			"result error(%s): %s", result.GetCode(), result.GetMessage())
	}
	return result, nil
}
