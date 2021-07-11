package types

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang/protobuf/proto"
	"github.com/paintmeyellow/tron-demo/api"
	"time"
)

func SignTxn(txn *api.TransactionExtention, priv *ecdsa.PrivateKey) error {
	txn.Transaction.RawData.Timestamp = time.Now().UnixNano() / 1000000

	rawData, err := proto.Marshal(txn.Transaction.RawData)
	if err != nil {
		return err
	}

	hash := sha256.Sum256(rawData)
	txn.Txid = hash[:]

	for range txn.Transaction.RawData.Contract {
		signature, err := crypto.Sign(hash[:], priv)
		if err != nil {
			return err
		}
		txn.Transaction.Signature = append(txn.Transaction.Signature, signature)
	}

	return nil
}
