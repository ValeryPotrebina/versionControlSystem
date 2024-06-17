package object

import (
	"bytes"
	"encoding/gob"
	"log"
)

// Tree elem that have children (dir for example)
type Tree struct {
	Children []Child
}

// Child of Tree.
type Child struct {
	Type uint   //TypeBlob or TypeTree.
	Name []byte //FileName or TreeName.
	Hash []byte //Hash of child object.
}

// Serialize tree
func (t *Tree) Serialize() (data []byte) {
	var b bytes.Buffer
	encoder := gob.NewEncoder(&b)
	err := encoder.Encode(t)
	if err != nil {
		log.Panic(err)
	}
	data = b.Bytes()
	return
}

// Deserialize data into tree
func DeserializeTree(data []byte) *Tree {
	var tree Tree
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&tree)
	if err != nil {
		log.Panic(err)
	}
	return &tree
}

// Create object for tree
func (t *Tree) CreateObject() *Object {
	return &Object{
		TypeTree,
		t.Serialize(),
	}
}

