package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"sort"

	"github.com/henkman/cryptopals"
)

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
		h += float64(cryptopals.HammingDistance(enc[:ks], enc[l-ks:]))
		h += float64(cryptopals.HammingDistance(enc[ks:], enc[:l-ks]))
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
					s := cryptopals.LanguageProbability(t, cryptopals.English_Letter_Frequency)
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
			if cryptopals.IsUnprintable(c) {
				fmt.Printf("0x%02x,", c)
			} else {
				fmt.Printf("%c,", c)
			}
		}
	}
}
