package main

import (
	"bytes"
	"encoding/base64"
	"io"
	"log"
	"os"

	"github.com/henkman/cryptopals"
)

func main() {
	key := bytes.NewReader([]byte("Terminator X: Bring the noise"))
	var data io.Reader
	{
		fd, err := os.Open("6.txt")
		if err != nil {
			log.Fatal(err)
		}
		defer fd.Close()
		data = base64.NewDecoder(base64.StdEncoding, fd)
	}
	if err := cryptopals.XorReader(data, key, os.Stdout); err != nil {
		log.Fatal(err)
	}
}
