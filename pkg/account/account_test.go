package account

import (
	"bytes"
	"crypto/ecdsa"
	"errors"
	"github.com/paintmeyellow/gotron/pkg/common/crypto"
	"testing"
)

func TestCreate(t *testing.T) {
	acc, err := Create()
	if err != nil {
		t.Fatal(err)
	}
	addr, err := crypto.Base58ToAddress(acc.Address)
	if err != nil {
		t.Fatal(err)
	}
	if !addr.Valid() {
		t.Fatal(err)
	}
	if ok, err := isAddresForPrivateKey(addr, acc.PrivateKey); err != nil || !ok {
		t.Fatalf("result:%t, error:%v", ok, err)
	}
}

func isAddresForPrivateKey(addr *crypto.Address, priv string) (bool, error) {
	privateKey, err := crypto.HexToECDSA(priv)
	if err != nil {
		return false, err
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return false, err
	}
	addrFromPriv := crypto.PubkeyToAddress(*publicKeyECDSA)
	if !bytes.Equal(addrFromPriv.Bytes(), addr.Bytes()) {
		return false, errors.New("invalid private key")
	}
	return true, nil
}
