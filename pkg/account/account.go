package account

import (
	"crypto/ecdsa"
	"encoding/hex"
	"github.com/paintmeyellow/tron-demo/pkg/common/base58"
	"github.com/pkg/errors"
)

var ErrAssertType = errors.New("cannot assert type")

type Account struct {
	Address    string
	PrivateKey string
}

func Create() (*Account, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.Wrap(ErrAssertType, "publicKey is not of type *ecdsa.PublicKey")
	}

	addresBytes := crypto.PubkeyToAddress(*publicKeyECDSA).Bytes()
	tronAddress := base58.EncodeCheck(addresBytes)

	return &Account{
		Address:    tronAddress,
		PrivateKey: hex.EncodeToString(crypto.FromECDSA(privateKey)),
	}, nil
}
