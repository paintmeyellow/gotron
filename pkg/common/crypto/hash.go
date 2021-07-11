package crypto

import "github.com/paintmeyellow/tron-demo/common/hexutil"

const (
	// HashLength is the expected length of the hash
	HashLength = 32
)

// Hash represents the 32 byte Keccak256 hash of arbitrary data.
type Hash [HashLength]byte

// BytesToHash sets b to hash.
// If b is larger than len(h), b will be cropped from the left.
func BytesToHash(b []byte) *Hash {
	var h Hash
	h.SetBytes(b)
	return &h
}

// HexToHash sets byte representation of s to hash.
// If b is larger than len(h), b will be cropped from the left.
func HexToHash(s string) (*Hash, error) {
	b, err := hexutil.Decode(s)
	if err != nil {
		return nil, err
	}
	return BytesToHash(b), nil
}

// Bytes gets the byte representation of the underlying hash.
func (h Hash) Bytes() []byte {
	return h[:]
}

// Hex converts a hash to a hex string.
func (h Hash) Hex() string {
	return hexutil.Bytes2Hex(h[:])
}

// SetBytes sets the hash to the value of b.
// If b is larger than len(h), b will be cropped from the left.
func (h *Hash) SetBytes(b []byte) {
	if len(b) > len(h) {
		b = b[len(b)-HashLength:]
	}
	copy(h[HashLength-len(b):], b)
}
