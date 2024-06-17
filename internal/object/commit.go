package object

import (
	"bytes"
	"encoding/gob"
	"log"
)

type Commit struct {
	Origin      []byte		//Reference to prev commit 
	Tree        []byte
	Author      []byte
	Time        int64
	Description []byte
}

func (c *Commit) Serialize() (data []byte) {
	var b bytes.Buffer
	encoder := gob.NewEncoder(&b)
	err := encoder.Encode(c)
	if err != nil {
		log.Panic(err)
	}
	data = b.Bytes()
	return
}

// Deserialize data into tree
func DeserializeCommit(data []byte) *Commit {
	var commit Commit
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&commit)
	if err != nil {
		log.Panic(err)
	}
	return &commit
}

func (c *Commit) CreateObject() *Object {
	return &Object{
		TypeCommit,
		c.Serialize(),
	}
}
