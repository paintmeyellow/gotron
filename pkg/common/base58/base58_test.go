package base58

import (
	"github.com/paintmeyellow/tron-demo/common/hexutil"
	"strings"
	"testing"
)

func TestEncodeCheck(t *testing.T) {
	addrBytes, err := hexutil.Decode("41c15f2021507b4f4b92bff5ff17f2f0c67ee6e7bf")
	if err != nil {
		t.Fatal(err)
	}
	encode := EncodeCheck(addrBytes)
	if encode != "TTbfNX2MbRJ1NG1FqQmKUy4E24sV9K2yzj" {
		t.Fatalf("failure: %s", encode)
	}
}

func TestDecodeCheck(t *testing.T) {
	decodeBytes, err := DecodeCheck("TTbfNX2MbRJ1NG1FqQmKUy4E24sV9K2yzj")
	if err != nil {
		t.Fatal(err)
	}
	decode := hexutil.Encode(decodeBytes)
	if !strings.EqualFold(decode, "41c15f2021507b4f4b92bff5ff17f2f0c67ee6e7bf") {
		t.Fatalf("failure: %s", decode)
	}
}
