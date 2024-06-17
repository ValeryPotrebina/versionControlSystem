package storage

import (
	"bytes"
	"encoding/gob"
	"log"
)

func SerializeRefs(refs map[string][]byte) []byte {
	var b bytes.Buffer
	encoder := gob.NewEncoder(&b)
	err := encoder.Encode(refs)
	if err != nil {
		log.Panic(err)
	}
	return b.Bytes()
}

func DeserializeRefs(data []byte) map[string][]byte {
	var m map[string][]byte
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&m)
	if err != nil {
		log.Panic(err)
	}
	return m
}
