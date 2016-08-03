package main

import (
	"bytes"
	"encoding/hex"
	"io"
	"log"
	"os"
)

func xor(data io.Reader, key io.ReadSeeker, out io.Writer) error {
	const BUF_SIZE = 32 * 1024

	// short path for key files
	// smaller than buffer
	if kf, ok := key.(*os.File); ok {
		fi, err := kf.Stat()
		if err != nil {
			return err
		}
		if fi.Size() < BUF_SIZE {
			skbuf := make([]byte, fi.Size())
			_, err := key.Read(skbuf)
			if err != nil {
				return err
			}
			key = bytes.NewReader(skbuf)
		}
	}

	var dbuf, kbuf [BUF_SIZE]byte
	for {
		n, err := data.Read(dbuf[:])
		if n > 0 {
			{
				o := 0
				for o != n {
					kn, err := key.Read(kbuf[o:n])
					o += kn
					if err != nil {
						if err == io.EOF {
							key.Seek(0, 0)
							continue
						}
						return err
					}
				}
			}
			for i := 0; i < n; i++ {
				dbuf[i] ^= kbuf[i]
			}
			_, err := out.Write(dbuf[:n])
			if err != nil {
				return err
			}
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
	}
	return nil
}

func main() {
	data, err := hex.DecodeString("1c0111001f010100061a024b53535009181c")
	if err != nil {
		log.Fatal(err)
	}
	key, err := hex.DecodeString("686974207468652062756c6c277320657965")
	if err != nil {
		log.Fatal(err)
	}
	if err := xor(bytes.NewBuffer(data), bytes.NewReader(key), os.Stdout); err != nil {
		log.Fatal(err)
	}
}
