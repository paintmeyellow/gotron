package trongrid

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type WalletClient struct {
	Client  *http.Client
	Address url.URL
}

func NewWalletClient(address url.URL) *WalletClient {
	return &WalletClient{
		Client:  http.DefaultClient,
		Address: address,
	}
}

type PaginatedTxns struct {
	Data []*Transaction `json:"data"`
	Meta struct {
		Fingerprint string `json:"fingerprint,omitempty"`
	} `json:"meta"`
}

type Transaction struct {
	TxID        string          `json:"txID"`
	RawData     *TransactionRaw `json:"raw_data,omitempty"`
	BlockNumber int64           `json:"blockNumber,omitempty"`
}

type TransactionRaw struct {
	Contract  []*TransactionContract `json:"contract,omitempty"`
	Timestamp int64                  `json:"timestamp,omitempty"`
}

type TransactionContract struct {
	Parameter struct {
		Value struct {
			Amount int64 `json:"amount"`
		} `json:"value"`
	} `json:"parameter,omitempty"`
}

func (tx *Transaction) Amount() int64 {
	if len(tx.RawData.Contract) == 0 {
		return 0
	}
	return tx.RawData.Contract[0].Parameter.Value.Amount
}

func (tx *Transaction) Timestamp() int64 {
	return tx.RawData.Timestamp
}

type Block struct {
	ID     string `json:"blockID"`
	Header struct {
		RawData struct {
			Number    int64 `json:"number"`
			Timestamp int64 `json:"timestamp"`
		} `json:"raw_data"`
	} `json:"block_header"`
}

func (b *Block) Number() int64 {
	return b.Header.RawData.Number
}

type AccTxnsContract struct {
	Address      string
	MinTimestamp int64
	OnlyTo       bool
	Limit        int64
}

type PaginatedAccTxnsContract struct {
	Address        string
	MinTimestamp   int64
	OnlyTo         bool
	SearchInternal bool
	Limit          int64
	Fingerprint    string
}

//GetAccountTransactionsList query the list of all normal transactions
func (c *WalletClient) GetAccountTransactionsList(ctx context.Context, contr AccTxnsContract) ([]*Transaction, error) {
	var (
		page        = 1
		fingerprint string
		list        = make([]*Transaction, 0)
	)

	for fingerprint != "" || page == 1 {
		batch, err := c.GetPaginatedAccountTransactions(ctx, PaginatedAccTxnsContract{
			Address:        contr.Address,
			MinTimestamp:   contr.MinTimestamp,
			OnlyTo:         contr.OnlyTo,
			SearchInternal: false,
			Limit:          contr.Limit,
			Fingerprint:    fingerprint,
		})
		if err != nil {
			return nil, err
		}
		list = append(list, batch.Data...)
		fingerprint = batch.Meta.Fingerprint
		page++
	}
	return list, nil
}

//GetPaginatedAccountTransactions query the paginated list of all normal transactions
func (c *WalletClient) GetPaginatedAccountTransactions(ctx context.Context, contr PaginatedAccTxnsContract) (*PaginatedTxns, error) {
	addr := c.Address
	addr.Path += fmt.Sprintf("/v1/accounts/%s/transactions", contr.Address)
	params := url.Values{
		"only_to":         {strconv.FormatBool(contr.OnlyTo)},
		"min_timestamp":   {strconv.FormatInt(contr.MinTimestamp, 10)},
		"search_internal": {strconv.FormatBool(contr.SearchInternal)},
		"limit":           {strconv.FormatInt(contr.Limit, 10)},
		"fingerprint":     {contr.Fingerprint},
	}
	addr.RawQuery = params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, addr.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var txns PaginatedTxns
	if err = json.NewDecoder(resp.Body).Decode(&txns); err != nil {
		return nil, err
	}
	return &txns, nil
}

func (c *WalletClient) GetNowBlock(ctx context.Context) (*Block, error) {
	addr := c.Address
	addr.Path += "/wallet/getnowblock"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, addr.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var block Block
	if err = json.NewDecoder(resp.Body).Decode(&block); err != nil {
		return nil, err
	}
	return &block, nil
}
