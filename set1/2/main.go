package main

import (
	"bytes"
	"encoding/hex"
	"log"
	"os"

	"github.com/henkman/cryptopals"
)

func main() {
	data, err := hex.DecodeString("1c0111001f010100061a024b53535009181c")
	if err != nil {
		log.Fatal(err)
	}
	key, err := hex.DecodeString("686974207468652062756c6c277320657965")
	if err != nil {
		log.Fatal(err)
	}
	if err := cryptopals.XorReader(bytes.NewBuffer(data), bytes.NewReader(key), os.Stdout); err != nil {
		log.Fatal(err)
	}
}
