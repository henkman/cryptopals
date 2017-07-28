package cryptopals

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"
)

func TestXor(t *testing.T) {
	src := []byte{0x2, 0x3, 0x1, 0x3, 0x1, 0x2}
	key := []byte{0x1, 0x2, 0x3}
	Xor(src, key, src)
	want := []byte{0x3, 0x1, 0x2, 0x2, 0x3, 0x1}
	if !bytes.Equal(src, want) {
		t.Fail()
		fmt.Println("have:", hex.EncodeToString(src), "want:", hex.EncodeToString(want))
	}
}
