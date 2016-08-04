package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
)

var (
	_file     string
	_out      string
	_csv      bool
	_go       bool
	_sortfreq bool
)

func init() {
	flag.StringVar(&_file, "f", "", "file to read")
	flag.StringVar(&_out, "o", "", "file that will contain frequencies")
	flag.BoolVar(&_csv, "csv", false, "print as csv")
	flag.BoolVar(&_go, "go", false, "print as go array")
	flag.BoolVar(&_sortfreq, "s", false, "sort by frequency, not by ASCII code")
	flag.Parse()
}

func isUnprintable(c byte) bool {
	return c != '\t' && c != '\n' && c != '\r' && (c >= 0 && c <= 0x1F) || c == 0x7F
}

type CharacterFrequency struct {
	Character byte
	Frequency float32
}

type ByFrequency []CharacterFrequency

func (a ByFrequency) Len() int           { return len(a) }
func (a ByFrequency) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByFrequency) Less(i, j int) bool { return a[i].Frequency > a[j].Frequency }

func main() {
	if _file == "" || _out == "" || _csv == _go {
		flag.Usage()
		return
	}
	var lang [128]CharacterFrequency
	for i, _ := range lang {
		lang[i].Character = byte(i)
	}
	{
		fd, err := os.Open(_file)
		if err != nil {
			log.Fatal(err)
		}
		defer fd.Close()
		var total float32
		var buf [1024 * 8]byte
		for {
			n, err := fd.Read(buf[:])
			if n > 0 {
				for _, c := range buf[:n] {
					if c > 0x7F {
						log.Fatal("file contains character out of ASCII range")
					} else if isUnprintable(c) {
						log.Fatal("file contains unprintable characters")
					}
					lang[c].Frequency++
				}
				total += float32(n)
			}
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Fatal(err)
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
		fd, err := os.OpenFile(_out, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0750)
		if err != nil {
			log.Fatal(err)
		}
		defer fd.Close()
		if _csv {
			fd.WriteString("character,frequency\n")
			for _, cf := range lang {
				if isUnprintable(cf.Character) {
					continue
				} else if cf.Character == ' ' {
					fmt.Fprintf(fd, "SPACE,%f\n", cf.Frequency)
				} else if cf.Character == '\t' {
					fmt.Fprintf(fd, "TAB,%f\n", cf.Frequency)
				} else if cf.Character == '\n' {
					fmt.Fprintf(fd, "NL,%f\n", cf.Frequency)
				} else if cf.Character == '\r' {
					fmt.Fprintf(fd, "CR,%f\n", cf.Frequency)
				} else {
					fmt.Fprintf(fd, "%c%f\n", cf.Character, cf.Frequency)
				}
			}
		} else if _go {
			fd.WriteString("var LANG = [128]byte {\n")
			for _, cf := range lang {
				if isUnprintable(cf.Character) {
					continue
				} else if cf.Character == ' ' {
					fmt.Fprintf(fd, "\t' ':%f,\n", cf.Frequency)
				} else if cf.Character == '\t' {
					fmt.Fprintf(fd, "\t'\\t':%f,\n", cf.Frequency)
				} else if cf.Character == '\n' {
					fmt.Fprintf(fd, "\t'\\n':%f,\n", cf.Frequency)
				} else if cf.Character == '\r' {
					fmt.Fprintf(fd, "\t'\\r':%f,\n", cf.Frequency)
				} else if cf.Character == '\'' {
					fmt.Fprintf(fd, "\t'\\'':%f,\n", cf.Frequency)
				} else {
					fmt.Fprintf(fd, "\t'%c':%f,\n", cf.Character, cf.Frequency)
				}
			}
			fd.WriteString("}\n")
		}
	}
}
