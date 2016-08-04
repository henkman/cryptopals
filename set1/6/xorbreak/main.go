package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"sort"
)

func hammingDistance(a, b []byte) uint {
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

func languageProbability(b []byte, lang [128]float64) float64 {
	isUnprintable := func(c byte) bool {
		return c != '\t' && c != '\n' && c != '\r' && c >= 0 && c <= 0x1F
	}
	cosineSimilarity := func(a []float64, b []float64) float64 {
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
			if k >= length_a {
				s2 += math.Pow(b[k], 2)
				continue
			}
			if k >= length_b {
				s1 += math.Pow(a[k], 2)
				continue
			}
			sumA += a[k] * b[k]
			s1 += math.Pow(a[k], 2)
			s2 += math.Pow(b[k], 2)
		}
		if s1 == 0 || s2 == 0 {
			return 0
		}
		return sumA / (math.Sqrt(s1) * math.Sqrt(s2))
	}
	var freq [128]float64
	for _, c := range b {
		if c > 0x7F || isUnprintable(c) {
			return 0
		}
		freq[c]++
	}
	return cosineSimilarity(freq[:], lang[:])
}

func englishProbability(b []byte) float64 {
	var FREQ = [128]float64{
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
	return languageProbability(b, FREQ)
}

type KeysizeProbability struct {
	Keysize     uint
	Probability float64
}
type KeysizeByProbability []KeysizeProbability

func (p KeysizeByProbability) Len() int {
	return len(p)
}
func (p KeysizeByProbability) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
func (p KeysizeByProbability) Less(i, j int) bool {
	a := p[i]
	b := p[j]
	if a.Probability == b.Probability {
		return a.Keysize < b.Keysize
	}
	return a.Probability < b.Probability
}

func XORFindProbableKeysizes(enc []byte, limit uint) []KeysizeProbability {
	l := len(enc)
	ksps := make([]KeysizeProbability, 0, l/2-2)
	for ks := 2; ks < l/2; ks++ {
		var h float64
		h += float64(hammingDistance(enc[:ks], enc[l-ks:]))
		h += float64(hammingDistance(enc[ks:], enc[:l-ks]))
		h /= float64(l)
		ksps = append(ksps, KeysizeProbability{uint(ks), h})
	}
	sort.Sort(KeysizeByProbability(ksps))
	if uint(len(ksps)) > limit {
		ksps = ksps[:limit]
	}
	return ksps
}

type CharacterProbability struct {
	Character   byte
	Probability float64
}

type CharacterByProbability []CharacterProbability

func (p CharacterByProbability) Len() int {
	return len(p)
}
func (p CharacterByProbability) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
func (p CharacterByProbability) Less(i, j int) bool {
	return p[i].Probability < p[j].Probability
}

type CharacterCount struct {
	Character byte
	Count     uint
}

type CharacterByCount []CharacterCount

func (p CharacterByCount) Len() int {
	return len(p)
}
func (p CharacterByCount) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
func (p CharacterByCount) Less(i, j int) bool {
	return p[i].Count > p[j].Count
}

func isUnprintable(c byte) bool {
	return c != '\t' && c != '\n' && c != '\r' && (c >= 0 && c <= 0x1F) || c == 0x7F
}

var (
	_file                    string
	_b64                     bool
	_printprobabilities      bool
	_englishprobabilitylimit float64
	_keysizelimit            uint
)

func init() {
	flag.StringVar(&_file, "f", "", "file to read")
	flag.BoolVar(&_b64, "b64", false, "file is base64 encoded")
	flag.BoolVar(&_printprobabilities, "pbs", false, "print all found probable characters")
	flag.Float64Var(&_englishprobabilitylimit, "epl", 0.7, "limit for the english probability of the single character phase")
	flag.UintVar(&_keysizelimit, "ksl", 20, "limits further calculations to the most probable keysizes")
	flag.Parse()
}

func main() {
	if _file == "" {
		flag.Usage()
		return
	}
	var enc []byte
	{
		tenc, err := ioutil.ReadFile(_file)
		if err != nil {
			log.Fatal(err)
		}
		if _b64 {
			b64, err := base64.StdEncoding.DecodeString(string(tenc))
			if err != nil {
				log.Fatal(err)
			}
			enc = b64
		} else {
			enc = tenc
		}
	}
	{
		ksps := XORFindProbableKeysizes(enc, _keysizelimit)
		var lks uint
		for _, h := range ksps {
			if h.Keysize > lks {
				lks = h.Keysize
			}
		}
		pcbs := make([][]CharacterProbability, lks)
		for _, h := range ksps {
			chunks := uint(len(enc)) / h.Keysize
			if uint(len(enc))%h.Keysize != 0 {
				chunks++
			}
			t := make([]byte, chunks)
			var p uint
			for p = 0; p < h.Keysize; p++ {
				cpbs := make([]CharacterProbability, 0, 256)
				for c := 1; c < 255; c++ {
					o := 0
					var i uint
					for i = p; i < uint(len(enc)); i += h.Keysize {
						t[o] = enc[i] ^ byte(c)
						o++
					}
					s := englishProbability(t)
					if s >= _englishprobabilitylimit {
						cpbs = append(cpbs, CharacterProbability{byte(c), s})
					}
				}
				if len(cpbs) == 0 {
					continue
				}
				if pcbs[p] == nil {
					pcbs[p] = append([]CharacterProbability{}, cpbs...)
				} else {
					pcbs[p] = append(pcbs[p], cpbs...)
				}
			}
		}
		mpk := make([]byte, 0, lks)
		cc := make([]CharacterCount, 256)
		for i := 0; i < 256; i++ {
			cc[i].Character = byte(i)
		}
		for p, cpbs := range pcbs {
			if len(cpbs) == 0 {
				continue
			}
			if _printprobabilities {
				fmt.Printf("p=%d:\n", p)
			}
			for i := 0; i < 256; i++ {
				cc[i].Count = 0
			}
			sort.Sort(CharacterByProbability(cpbs))
			for _, cpb := range cpbs {
				cc[cpb.Character].Count++
				if _printprobabilities {
					fmt.Printf("char=%c, prob=%f\n", cpb.Character, cpb.Probability)
				}
			}
			ok := false
			mkci := 0
			for i := 1; i < len(cc); i++ {
				if cc[i].Count > cc[mkci].Count {
					mkci = i
					ok = true
				}
			}
			if ok {
				mpk = append(mpk, cc[mkci].Character)
			}
		}
		fmt.Println("most probable key:")
		for _, c := range mpk {
			if isUnprintable(c) {
				fmt.Printf("0x%02x,", c)
			} else {
				fmt.Printf("%c,", c)
			}
		}
	}
}
