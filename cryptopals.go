package cryptopals

import (
	"bytes"
	"crypto/aes"
	"io"
	"math"
	"os"
)

var (
	English_Letter_Frequency = [128]float64{
		' ':  0.169517,
		'e':  0.096241,
		't':  0.070165,
		'a':  0.062427,
		'o':  0.059632,
		'n':  0.054496,
		'h':  0.049977,
		'i':  0.049460,
		's':  0.048570,
		'r':  0.043618,
		'd':  0.034301,
		'l':  0.031643,
		'u':  0.022752,
		'm':  0.018603,
		'w':  0.018392,
		'c':  0.017356,
		'f':  0.016327,
		'g':  0.016189,
		'y':  0.015914,
		',':  0.015488,
		'p':  0.012140,
		'b':  0.011687,
		'.':  0.008415,
		'v':  0.007032,
		'k':  0.006669,
		'"':  0.005339,
		'I':  0.004507,
		'\'': 0.003816,
		'-':  0.003653,
		';':  0.002284,
		'T':  0.001912,
		'A':  0.001464,
		'M':  0.001357,
		'S':  0.001324,
		'H':  0.001309,
		'!':  0.001174,
		'W':  0.001172,
		'B':  0.001056,
		'?':  0.001056,
		'x':  0.000915,
		'q':  0.000831,
		'C':  0.000767,
		'j':  0.000732,
		'L':  0.000725,
		'D':  0.000718,
		'_':  0.000691,
		'E':  0.000652,
		'N':  0.000589,
		'z':  0.000559,
		'P':  0.000524,
		'O':  0.000493,
		'Y':  0.000488,
		'F':  0.000411,
		'J':  0.000385,
		'G':  0.000360,
		'R':  0.000348,
		':':  0.000337,
		'K':  0.000121,
		'Q':  0.000120,
		')':  0.000112,
		'(':  0.000112,
		'V':  0.000100,
		'U':  0.000099,
		'0':  0.000056,
		'1':  0.000055,
		'*':  0.000051,
		'X':  0.000032,
		'2':  0.000029,
		'8':  0.000024,
		'5':  0.000023,
		'7':  0.000020,
		'3':  0.000020,
		'4':  0.000019,
		'6':  0.000015,
		'9':  0.000013,
		'Z':  0.000011,
		'&':  0.000002,
		'[':  0.000001,
		']':  0.000001,
		'$':  0.000001,
		'/':  0.000000,
		'>':  0.000000,
		'%':  0.000000,
		'#':  0.000000,
		'@':  0.000000,
		'+':  0.000000,
		'<':  0.000000,
		'\\': 0.000000,
		'=':  0.000000,
		'^':  0.000000,
		'`':  0.000000,
		'\r': 0.000000,
		'\n': 0.000000,
		'\t': 0.000000,
		'{':  0.000000,
		'|':  0.000000,
		'}':  0.000000,
		'~':  0.000000,
	}
)

func Xor(src, key, dst []byte) {
	l := len(key)
	for i := range src {
		dst[i] = src[i] ^ key[i%l]
	}
}

func XorReader(src io.Reader, key io.ReadSeeker, dst io.Writer) error {
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
		n, err := src.Read(dbuf[:])
		if n > 0 {
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
			for i := 0; i < n; i++ {
				dbuf[i] ^= kbuf[i]
			}
			_, err := dst.Write(dbuf[:n])
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

func Pkcs7PaddingBytesNeeded(src []byte, blocksize uint) uint {
	l := uint(len(src))
	r := l % blocksize
	if r == 0 {
		return 0
	}
	return blocksize - r
}

func EncryptAESCBC(src, dst, iv, key []byte) error {
	c, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	bs := c.BlockSize()
	for i := 0; i < len(src)/bs; i++ {
		a := src[i*bs:]
		b := dst[i*bs:]
		Xor(iv, a[:bs], iv)
		c.Encrypt(b, iv)
		copy(iv, b)
	}
	return nil
}

func DecryptAESCBC(src, dst, iv, key []byte) error {
	c, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	bs := c.BlockSize()
	for i := 0; i < len(src)/bs; i++ {
		a := src[i*bs:]
		b := dst[i*bs:]
		c.Decrypt(b, a)
		Xor(b[:bs], iv, b[:bs])
		copy(iv, a)
	}
	return nil
}

func HammingDistance(a, b []byte) uint {
	var TABLE = [256]byte{
		0, 1, 1, 2, 1, 2, 2, 3, 1, 2, 2, 3, 2, 3, 3, 4,
		1, 2, 2, 3, 2, 3, 3, 4, 2, 3, 3, 4, 3, 4, 4, 5,
		1, 2, 2, 3, 2, 3, 3, 4, 2, 3, 3, 4, 3, 4, 4, 5,
		2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6,
		1, 2, 2, 3, 2, 3, 3, 4, 2, 3, 3, 4, 3, 4, 4, 5,
		2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6,
		2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6,
		3, 4, 4, 5, 4, 5, 5, 6, 4, 5, 5, 6, 5, 6, 6, 7,
		1, 2, 2, 3, 2, 3, 3, 4, 2, 3, 3, 4, 3, 4, 4, 5,
		2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6,
		2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6,
		3, 4, 4, 5, 4, 5, 5, 6, 4, 5, 5, 6, 5, 6, 6, 7,
		2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6,
		3, 4, 4, 5, 4, 5, 5, 6, 4, 5, 5, 6, 5, 6, 6, 7,
		3, 4, 4, 5, 4, 5, 5, 6, 4, 5, 5, 6, 5, 6, 6, 7,
		4, 5, 5, 6, 5, 6, 6, 7, 5, 6, 6, 7, 6, 7, 7, 8,
	}
	var n int
	if len(a) >= len(b) {
		n = len(b)
	} else {
		n = len(a)
	}
	var h uint
	for i := 0; i < n; i++ {
		h += uint(TABLE[a[i]^b[i]])
	}
	return h
}

func IsUnprintable(c byte) bool {
	return c != '\t' && c != '\n' && c != '\r' && (c >= 0 && c <= 0x1F) || c == 0x7F
}

func CosineSimilarity(a []float64, b []float64) float64 {
	cn := 0
	lena := len(a)
	lenb := len(b)
	if lena > lenb {
		cn = lena
	} else {
		cn = lenb
	}
	suma := 0.0
	s1 := 0.0
	s2 := 0.0
	for k := 0; k < cn; k++ {
		if k >= lena {
			s2 += b[k] * b[k]
			continue
		}
		if k >= lenb {
			s1 += a[k] * a[k]
			continue
		}
		suma += a[k] * b[k]
		s1 += a[k] * a[k]
		s2 += b[k] * b[k]
	}
	if s1 == 0 || s2 == 0 {
		return 0
	}
	return suma / (math.Sqrt(s1) * math.Sqrt(s2))
}

func LanguageProbability(b []byte, lang [128]float64) float64 {
	var freq [128]float64
	for _, c := range b {
		if c > 0x7F || IsUnprintable(c) {
			return 0
		}
		freq[c]++
	}
	return CosineSimilarity(freq[:], lang[:])
}
