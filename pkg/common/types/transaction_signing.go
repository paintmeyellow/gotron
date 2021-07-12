package types

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/paintmeyellow/gotron/pkg/proto/core"
	"google.golang.org/protobuf/proto"
	"time"
)

func SignTxn(txn *core.Transaction, priv *ecdsa.PrivateKey) ([]byte, error) {
	txn.RawData.Timestamp = time.Now().UnixNano() / 1000000

	rawData, err := proto.Marshal(txn.RawData)
	if err != nil {
		return nil, err
	}
	hash := sha256.Sum256(rawData)

	for range txn.RawData.Contract {
		signature, err := crypto.Sign(hash[:], priv)
		if err != nil {
			return nil, err
		}
		txn.Signature = append(txn.Signature, signature)
	}

	return hash[:], nil
}
