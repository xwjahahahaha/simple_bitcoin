package blc

import (
	"testing"
)

func TestBase58(t *testing.T) {
	hash := []byte("1QrBDJhUzahDLk3aLRHauHH8NdNBWyihXT")
	t.Logf("hash = %x\n", hash)
	encode := Base58Encode(hash)
	//t.Logf("encode = %x\n", encode)
	t.Logf("decode = %s\n", Base58Decode(encode))
}