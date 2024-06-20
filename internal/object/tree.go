package object

import (
	"bytes"
	"encoding/gob"
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
func (t *Tree) Serialize() (data []byte, err error) {
	var b bytes.Buffer
	encoder := gob.NewEncoder(&b)
	err = encoder.Encode(t)
	data = b.Bytes()
	return
}

// Deserialize data into tree
func DeserializeTree(data []byte) (*Tree, error) {
	var tree Tree
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&tree)
	return &tree, err
}

// Create object for tree
func (t *Tree) CreateObject() (*Object, error) {
	data, err := t.Serialize()
	if err != nil {
		return nil, err
	}
	return &Object{
		TypeTree,
		data,
	}, nil
}
