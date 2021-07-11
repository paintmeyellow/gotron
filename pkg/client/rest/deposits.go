package rest

import (
	"context"
	"github.com/paintmeyellow/gotron/pkg/trongrid"
	"golang.org/x/sync/errgroup"
	"sync"
)

const listDepositsBatchSize = 20

type Deposit struct {
	Hash                 string
	Address              string
	Amount               int64
	Confirmations        int64
	TransactionTimestamp int64
}

func (c *Client) ListDeposits(ctx context.Context, minTimestamp int64, addresses []string) ([]*Deposit, error) {
	var (
		err          error
		wg           sync.WaitGroup
		currentBlock *trongrid.Block
		list         = make([]*Deposit, 0)
	)
	errs, ctx := errgroup.WithContext(ctx)

	wg.Add(1)
	errs.Go(func() error {
		defer wg.Done()
		currentBlock, err = c.WalletClient.GetNowBlock(ctx)
		return err
	})

	for _, addr := range addresses {
		addr := addr
		errs.Go(func() error {
			txns, err := c.WalletClient.GetAccountTransactionsList(ctx, trongrid.AccTxnsContract{
				Address:      addr,
				MinTimestamp: minTimestamp,
				OnlyTo:       true,
				Limit:        listDepositsBatchSize,
			})
			if err != nil {
				return err
			}
			wg.Wait()
			for _, txn := range txns {
				confirmations := currentBlock.Number() - txn.BlockNumber + 1
				list = append(list, &Deposit{
					Hash:                 txn.TxID,
					Amount:               txn.Amount(),
					Address:              addr,
					Confirmations:        confirmations,
					TransactionTimestamp: txn.Timestamp(),
				})
			}
			return nil
		})
	}

	return list, errs.Wait()
}
