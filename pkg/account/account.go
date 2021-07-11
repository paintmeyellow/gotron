package account

import (
	"crypto/ecdsa"
	"encoding/hex"
	"github.com/paintmeyellow/gotron/pkg/common/base58"
	"github.com/paintmeyellow/gotron/pkg/common/crypto"
	"github.com/pkg/errors"
)

var ErrCreatePublicKey = errors.New("cannot create public key")

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
		return nil, ErrCreatePublicKey
	}

	addresBytes := crypto.PubkeyToAddress(*publicKeyECDSA).Bytes()
	tronAddress := base58.EncodeCheck(addresBytes)

	return &Account{
		Address:    tronAddress,
		PrivateKey: hex.EncodeToString(crypto.FromECDSA(privateKey)),
	}, nil
}
