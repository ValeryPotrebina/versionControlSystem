package object

import (
	"bytes"
	"compress/zlib"
	"crypto/sha256"
	"io"
	"log"
)

// Calculate sha256 hash of data
func CalculateHash(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

// Zip data with zlib
func Zip(data []byte) []byte {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	_, err := w.Write(data)
	if err != nil {
		log.Panic(err)
	}
	err = w.Close()
	if err != nil {
		log.Panic(err)
	}
	return b.Bytes()
}

// Unzip data with zlib
func Unzip(data []byte) []byte {
	var res bytes.Buffer
	r, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		log.Panic(err)
	}
	_, err = io.Copy(&res, r)
	if err != nil {
		log.Panic(err)
	}
	err = r.Close()
	if err != nil {
		log.Panic(err)
	}
	return res.Bytes()
}
