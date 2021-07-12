package types

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/paintmeyellow/gotron/pkg/proto/api"
	"google.golang.org/protobuf/proto"
	"time"
)

func SignTxn(txn *api.TransactionExtention, priv *ecdsa.PrivateKey) ([]byte, error) {
	txn.Transaction.RawData.Timestamp = time.Now().UnixNano() / 1000000

	rawData, err := proto.Marshal(txn.Transaction.RawData)
	if err != nil {
		return nil, err
	}

	hash := sha256.Sum256(rawData)

	for range txn.Transaction.RawData.Contract {
		signature, err := crypto.Sign(hash[:], priv)
		if err != nil {
			return nil, err
		}
		txn.Transaction.Signature = append(txn.Transaction.Signature, signature)
	}

	return hash[:], nil
}
