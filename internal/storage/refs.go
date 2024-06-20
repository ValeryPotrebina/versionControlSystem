package storage

import (
	"bytes"
	"encoding/gob"
)

func SerializeRefs(refs map[string][]byte) ([]byte, error) {
	var b bytes.Buffer
	encoder := gob.NewEncoder(&b)
	err := encoder.Encode(refs)
	return b.Bytes(), err
}

func DeserializeRefs(data []byte) (map[string][]byte, error) {
	var m map[string][]byte
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&m)
	return m, err
}
