package object

import (
	"bytes"
	"compress/zlib"
	"crypto/sha256"
	"io"
)

// Calculate sha256 hash of data
func CalculateHash(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

// Zip data with zlib
func Zip(data []byte) ([]byte, error) {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	_, err := w.Write(data)
	if err != nil {
		return nil, err
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// Unzip data with zlib
func Unzip(data []byte) ([]byte, error) {
	var res bytes.Buffer
	r, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(&res, r)
	if err != nil {
		return nil, err
	}
	err = r.Close()
	if err != nil {
		return nil, err
	}
	return res.Bytes(), nil
}
