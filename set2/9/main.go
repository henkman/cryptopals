package main

import (
	"encoding/hex"
	"fmt"

	"github.com/henkman/cryptopals"
)

func main() {
	const BS = 20
	unpadded := []byte("YELLOW SUBMARINE")
	l := cryptopals.Pkcs7PaddingBytesNeeded(unpadded, BS)
	fmt.Print(hex.EncodeToString(unpadded))
	var i uint
	for i = 0; i < l; i++ {
		fmt.Printf("%02x", byte(l))
	}
	fmt.Println()
}
