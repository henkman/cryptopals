package main

import (
	"encoding/hex"
	"fmt"
	"log"

	"github.com/henkman/cryptopals"
)

func main() {
	data, err := hex.DecodeString("1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b3736")
	if err != nil {
		log.Fatal(err)
	}
	dst := make([]byte, len(data))
	var c byte
	for ; c < 0xFF; c++ {
		cryptopals.Xor(data, []byte{c}, dst)
		s := cryptopals.LanguageProbability(dst, cryptopals.English_Letter_Frequency)
		if s >= 0.75 {
			fmt.Printf("'%s', %f key=%d\n", string(dst), s, c)
		}
	}
}
