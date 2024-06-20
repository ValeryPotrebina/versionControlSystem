package object

import (
	"bytes"
	"encoding/gob"
)

type Commit struct {
	Origin      []byte //Reference to prev commit
	Tree        []byte
	Author      []byte
	Time        int64
	Description []byte
}

func (c *Commit) Serialize() (data []byte, err error) {
	var b bytes.Buffer
	encoder := gob.NewEncoder(&b)
	err = encoder.Encode(c)
	data = b.Bytes()
	return
}

// Deserialize data into tree
func DeserializeCommit(data []byte) (*Commit, error) {
	var commit Commit
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&commit)
	return &commit, err
}

func (c *Commit) CreateObject() (*Object, error) {
	data, err := c.Serialize()
	if err != nil {
		return nil, err
	}
	return &Object{
		TypeCommit,
		data,
	}, nil
}
