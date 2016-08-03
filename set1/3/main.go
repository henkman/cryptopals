package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
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

func isLowercaseAlpha(c byte) bool {
	return c >= 'a' && c <= 'z'
}

func isUppercaseAlpha(c byte) bool {
	return c >= 'A' && c <= 'Z'
}

func isAlpha(c byte) bool {
	return isLowercaseAlpha(c) || isUppercaseAlpha(c)
}

func isNumeric(c byte) bool {
	return c >= '0' && c <= '9'
}

func isSymbol(c byte) bool {
	return (c >= '!' && c <= '/') ||
		(c >= ':' && c <= '@') ||
		(c >= '[' && c <= '`') ||
		(c >= '{' && c <= '~')
}

func isWhitespace(c byte) bool {
	return c == ' ' || c == '\n' || c == '\r' || c == '\t'
}

func isPrintable(c byte) bool {
	return isWhitespace(c) || isSymbol(c) || isNumeric(c) || isAlpha(c)
}

func englishScore(b []byte) float32 {
	hasSpace := false
	symbols := 0
	for _, c := range b {
		if !isPrintable(c) {
			return 0
		}
		if c == ' ' {
			hasSpace = true
		} else if isSymbol(c) {
			symbols++
		}
	}
	var score float32 = 0.5
	if hasSpace {
		score += 0.2
	}
	score -= (float32(symbols) / float32(len(b)))
	return score
}

func main() {
	data, err := hex.DecodeString("1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b3736")
	if err != nil {
		log.Fatal(err)
	}
	var c byte
	for ; c < 0xFF; c++ {
		key := bytes.NewReader([]byte{c})
		out := bytes.NewBufferString("")
		if err := xor(bytes.NewBuffer(data), key, out); err != nil {
			log.Fatal(err)
		}
		b := out.Bytes()
		if englishScore(b) >= 0.5 {
			fmt.Printf("'%s' key=%d\n", string(b), c)
		}
	}
}
