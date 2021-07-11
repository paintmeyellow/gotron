package crypto

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/crypto"
)

// GenerateKey generates a new private key.
func GenerateKey() (*ecdsa.PrivateKey, error) {
	return crypto.GenerateKey()
}

//FromECDSAPub gets public key bytes
func FromECDSAPub(pub *ecdsa.PublicKey) []byte {
	return crypto.FromECDSAPub(pub)
}

// FromECDSA exports a private key into a binary dump.
func FromECDSA(priv *ecdsa.PrivateKey) []byte {
	return crypto.FromECDSA(priv)
}

// HexToECDSA parses a secp256k1 private key.
func HexToECDSA(hexkey string) (*ecdsa.PrivateKey, error) {
	return crypto.HexToECDSA(hexkey)
}

func Sign(hash []byte, privateKey *ecdsa.PrivateKey) ([]byte, error) {
	return crypto.Sign(hash, privateKey)
}