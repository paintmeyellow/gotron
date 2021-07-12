package crypto

import (
	"crypto/ecdsa"
	"testing"
)

const (
	testPrivateKey = "001d58091ffe501203b7ff07dcd3a49d55560f93769a3eab0e7542741fe06555"
	testB58Address = "TFXEziu1kymNVHyxD5Ukb839ubSh1CZWfV"
	testHexAddress = "413ce79ab890eadeb5c38c15046a0e2e91846b9b65"
)

func TestBase58ToAddress(t *testing.T) {
	addr, err := Base58ToAddress(testB58Address)
	if err != nil {
		t.Fatal(err)
	}
	if testB58Address != addr.String() {
		t.Fatal("addresses are not the same")
	}
}

func TestHexToAddress(t *testing.T) {
	addr := HexToAddress(testHexAddress)
	if testB58Address != addr.String() || testHexAddress != addr.Hex() {
		t.Fatal("addresses are not the same")
	}
}

func TestPubkeyToAddress(t *testing.T) {
	privateKey, err := HexToECDSA(testPrivateKey)
	if err != nil {
		t.Fatal(err)
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		t.Fatal("publicKey is not type of *ecdsa.PublicKey")
	}
	addr := PubkeyToAddress(*publicKeyECDSA)
	t.Log(addr.String())
	if testB58Address != addr.String() {
		t.Fatal("addresses are not the same")
	}
}