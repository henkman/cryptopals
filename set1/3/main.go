package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"math"
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

func languageProbability(b []byte, lang [128]float32) float32 {
	isUnprintable := func(c byte) bool {
		return c != '\t' && c != '\n' && c != '\r' && c >= 0 && c <= 0x1F
	}
	cosineSimilarity := func(a []float32, b []float32) float32 {
		count := 0
		length_a := len(a)
		length_b := len(b)
		if length_a > length_b {
			count = length_a
		} else {
			count = length_b
		}
		sumA := 0.0
		s1 := 0.0
		s2 := 0.0
		for k := 0; k < count; k++ {
			bk64 := float64(b[k])
			ak64 := float64(a[k])
			if k >= length_a {
				s2 += math.Pow(bk64, 2)
				continue
			}
			if k >= length_b {
				s1 += math.Pow(ak64, 2)
				continue
			}
			sumA += ak64 * bk64
			s1 += math.Pow(ak64, 2)
			s2 += math.Pow(bk64, 2)
		}
		if s1 == 0 || s2 == 0 {
			return 0
		}
		return float32(sumA / (math.Sqrt(s1) * math.Sqrt(s2)))
	}
	var freq [128]float32
	for _, c := range b {
		if isUnprintable(c) {
			return 0
		}
		freq[c]++
	}
	return cosineSimilarity(freq[:], lang[:])
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
		s := languageProbability(b, FREQ)
		if s >= 0.6 {
			fmt.Printf("'%s', %f key=%d\n", string(b), s, c)
		}
	}
}
