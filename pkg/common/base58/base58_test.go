package base58

import (
	"github.com/paintmeyellow/gotron/pkg/common/hexutil"
	"strings"
	"testing"
)

const (
	testB58Address = "TFXEziu1kymNVHyxD5Ukb839ubSh1CZWfV"
	testHexAddress = "413ce79ab890eadeb5c38c15046a0e2e91846b9b65"
)

func TestEncodeCheck(t *testing.T) {
	addrBytes, err := hexutil.Decode(testHexAddress)
	if err != nil {
		t.Fatal(err)
	}
	encode := EncodeCheck(addrBytes)
	if encode != testB58Address {
		t.Fatalf("failure: %s", encode)
	}
}

func TestDecodeCheck(t *testing.T) {
	decodeBytes, err := DecodeCheck(testB58Address)
	if err != nil {
		t.Fatal(err)
	}
	decode := hexutil.Encode(decodeBytes)
	if !strings.EqualFold(decode, testHexAddress) {
		t.Fatalf("failure: %s", decode)
	}
}
