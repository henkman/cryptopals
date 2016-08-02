package main

import (
	"encoding/base64"
	"encoding/hex"
	"log"
	"os"
)

func main() {
	b, err := hex.DecodeString("49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d")
	if err != nil {
		log.Fatal(err)
	}
	w := base64.NewEncoder(base64.RawStdEncoding, os.Stdout)
	w.Write(b)
	w.Close()
}
