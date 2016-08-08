package main

import (
	"encoding/hex"
	"fmt"
)

func pkcs7PaddingNeeded(src []byte, blocksize uint) uint {
	l := uint(len(src))
	r := l % blocksize
	if r == 0 {
		return 0
	}
	return blocksize - r
}

func main() {
	const BS = 20
	unpadded := []byte("YELLOW SUBMARINE")
	l := pkcs7PaddingNeeded(unpadded, BS)
	fmt.Print(hex.EncodeToString(unpadded))
	var i uint
	for i = 0; i < l; i++ {
		fmt.Printf("%02x", byte(l))
	}
	fmt.Println()
}
