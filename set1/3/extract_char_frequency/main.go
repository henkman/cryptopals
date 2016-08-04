package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
)

var (
	_files    string
	_out      string
	_csv      bool
	_go       bool
	_sortfreq bool
	_skipnl   bool
)

func init() {
	flag.StringVar(&_files, "f", "", "commaseparated files to read")
	flag.StringVar(&_out, "o", "", "file that will contain frequencies")
	flag.BoolVar(&_csv, "csv", false, "print as csv")
	flag.BoolVar(&_go, "go", false, "print as go array")
	flag.BoolVar(&_sortfreq, "s", false, "sort by frequency, not by ASCII code")
	flag.BoolVar(&_skipnl, "snl", false, "skip newlines")
	flag.Parse()
}

func isUnprintable(c byte) bool {
	return c != '\t' && c != '\n' && c != '\r' && (c >= 0 && c <= 0x1F) || c == 0x7F
}

type CharacterFrequency struct {
	Character byte
	Frequency float64
}

type ByFrequency []CharacterFrequency

func (a ByFrequency) Len() int           { return len(a) }
func (a ByFrequency) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByFrequency) Less(i, j int) bool { return a[i].Frequency > a[j].Frequency }

func main() {
	if _files == "" || _csv == _go {
		flag.Usage()
		return
	}
	var lang [128]CharacterFrequency
	for i, _ := range lang {
		lang[i].Character = byte(i)
	}
	{
		var buf [1024 * 8]byte
		var total float64
		files := strings.Split(_files, ",")
		for _, file := range files {
			fd, err := os.Open(strings.TrimSpace(file))
			if err != nil {
				log.Fatal(err)
			}
			defer fd.Close()
			for {
				n, err := fd.Read(buf[:])
				if n > 0 {
					o := 0
					for _, c := range buf[:n] {
						if c > 0x7F || isUnprintable(c) {
							continue
						}
						if _skipnl && (c == '\n' || c == '\r') {
							continue
						}
						lang[c].Frequency++
						o++
					}
					total += float64(o)
				}
				if err != nil {
					if err == io.EOF {
						break
					}
					log.Fatal(err)
				}
			}
		}
		for i, _ := range lang {
			if lang[i].Frequency != 0 {
				lang[i].Frequency = lang[i].Frequency / total
			}
		}
	}
	if _sortfreq {
		sort.Sort(ByFrequency(lang[:]))
	}
	{
		var out io.WriteCloser
		if _out == "" {
			out = os.Stdout
		} else {
			fd, err := os.OpenFile(_out, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0750)
			if err != nil {
				log.Fatal(err)
			}
			defer fd.Close()
			out = fd
		}
		if _csv {
			io.WriteString(out, "character,frequency\n")
			for _, cf := range lang {
				if isUnprintable(cf.Character) {
					continue
				} else if cf.Character == ' ' {
					fmt.Fprintf(out, "SPACE,%f\n", cf.Frequency)
				} else if cf.Character == '\t' {
					fmt.Fprintf(out, "TAB,%f\n", cf.Frequency)
				} else if cf.Character == '\n' {
					fmt.Fprintf(out, "NL,%f\n", cf.Frequency)
				} else if cf.Character == '\r' {
					fmt.Fprintf(out, "CR,%f\n", cf.Frequency)
				} else {
					fmt.Fprintf(out, "%c%f\n", cf.Character, cf.Frequency)
				}
			}
		} else if _go {
			io.WriteString(out, "var FREQ = [128]float64 {\n")
			for _, cf := range lang {
				if isUnprintable(cf.Character) {
					continue
				} else if cf.Character == ' ' {
					fmt.Fprintf(out, "\t' ':%f,\n", cf.Frequency)
				} else if cf.Character == '\t' {
					fmt.Fprintf(out, "\t'\\t':%f,\n", cf.Frequency)
				} else if cf.Character == '\n' {
					fmt.Fprintf(out, "\t'\\n':%f,\n", cf.Frequency)
				} else if cf.Character == '\r' {
					fmt.Fprintf(out, "\t'\\r':%f,\n", cf.Frequency)
				} else if cf.Character == '\'' {
					fmt.Fprintf(out, "\t'\\'':%f,\n", cf.Frequency)
				} else if cf.Character == '\\' {
					fmt.Fprintf(out, "\t'\\\\':%f,\n", cf.Frequency)
				} else {
					fmt.Fprintf(out, "\t'%c':%f,\n", cf.Character, cf.Frequency)
				}
			}
			io.WriteString(out, "}\n")
		}
	}
}
