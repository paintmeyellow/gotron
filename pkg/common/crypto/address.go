package crypto

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/paintmeyellow/gotron/pkg/common/base58"
	"github.com/paintmeyellow/gotron/pkg/common/hexutil"
	"math/big"
)

const (
	// AddressLength is the expected length of the address
	AddressLength = 21
	// TronBytePrefix is the hex prefix to address
	TronBytePrefix = byte(0x41)
)

type Address [AddressLength]byte

func (a *Address) Bytes() []byte {
	return a[:]
}

func (a *Address) Hex() string {
	return hexutil.Bytes2Hex(a[:])
}

// String implements fmt.Stringer.
func (a Address) String() string {
	if a[0] == 0 {
		return new(big.Int).SetBytes(a.Bytes()).String()
	}
	return base58.EncodeCheck(a.Bytes())
}

func (a *Address) SetBytes(b []byte) {
	if len(b) > len(a) {
		b = b[len(b)-AddressLength:]
	}
	copy(a[AddressLength-len(b):], b)
}

func BytesToAddress(b []byte) *Address {
	var a Address
	a.SetBytes(b)
	return &a
}

// HexToAddress returns Address with byte values of s.
// If s is larger than len(h), s will be cropped from the left.
func HexToAddress(s string) *Address {
	addr, err := hexutil.Decode(s)
	if err != nil {
		return nil
	}
	return BytesToAddress(addr)
}

func Base58ToAddress(s string) (*Address, error) {
	addr, err := base58.DecodeCheck(s)
	if err != nil {
		return nil, err
	}
	return BytesToAddress(addr), nil
}

func PubkeyToAddress(pub ecdsa.PublicKey) *Address {
	address := crypto.PubkeyToAddress(pub)
	addressTron := make([]byte, 0)
	addressTron = append(addressTron, TronBytePrefix)
	addressTron = append(addressTron, address.Bytes()...)
	return BytesToAddress(addressTron)
}
