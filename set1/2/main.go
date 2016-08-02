package main

import (
	"bytes"
	"encoding/hex"
	"io"
	"log"
	"os"
)

func keyget(key io.ReadSeeker, buf []byte, toread int) error {
	o := 0
	for o != toread {
		n, err := key.Read(buf[o:toread])
		o += n
		if err != nil {
			if err == io.EOF {
				key.Seek(0, 0)
				continue
			}
			return err
		}
	}
	return nil
}

func xor(data io.Reader, key io.ReadSeeker, out io.Writer) error {
	dbuf := make([]byte, 32*1024)
	kbuf := make([]byte, len(dbuf))

	// short path for key files
	// smaller than buffer
	if kf, ok := key.(*os.File); ok {
		fi, err := kf.Stat()
		if err != nil {
			return err
		}
		if fi.Size() < int64(len(dbuf)) {
			skbuf := make([]byte, fi.Size())
			_, err := key.Read(skbuf)
			if err != nil {
				return err
			}
			key = bytes.NewReader(skbuf)
		}
	}

	for {
		n, err := data.Read(dbuf)
		if n > 0 {
			err = keyget(key, kbuf, n)
			if err != nil {
				return err
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
