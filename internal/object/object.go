package object

import (
	"bytes"
	"encoding/gob"
	"errors"
	"log"
)

// The object type.
const (
	TypeBlob = iota
	TypeTree
	TypeCommit
)

// Tree or Blob with field Type (TypeBlob or TypeTree). Data contains serialized Tree or Blob object
type Object struct {
	Type uint
	Data []byte
}

func TypeToString(t uint) string {
	switch t {
	case TypeBlob:
		return "Blob"
	case TypeTree:
		return "Tree"
	case TypeCommit:
		return "Commit"
	default:
		return ""
	}
}

// Calculate hash of the serialized object
func (o *Object) GetHash() []byte {
	return CalculateHash(o.Serialize())
}

// Get pair of hash and zipped data of object
func (o *Object) GetData() (hash []byte, data []byte) {
	b := o.Serialize()
	hash = CalculateHash(b)
	data = Zip(b)
	return
}

// Serialize object
func (o *Object) Serialize() []byte {
	var b bytes.Buffer
	encoder := gob.NewEncoder(&b)
	err := encoder.Encode(o)
	if err != nil {
		log.Panic(err)
	}
	return b.Bytes()
}

// Deserialize data into object
func DeserializeObject(data []byte) *Object {
	var obj Object
	decoder := gob.NewDecoder(bytes.NewReader(Unzip(data)))
	err := decoder.Decode(&obj)
	if err != nil {
		log.Panic(err)
	}
	return &obj
}

// Unpack object into tree if it is possible
func (o *Object) ParseTree() (*Tree, error) {
	if o.Type != TypeTree {
		return nil, errors.New("Object is not tree")
	}
	return DeserializeTree(o.Data), nil
}

// Unpack object into blob if it is possible
func (o *Object) ParseBlob() (*Blob, error) {
	if o.Type != TypeBlob {
		return nil, errors.New("Object is not blob")
	}
	return &Blob{o.Data}, nil
}

// Unpack object into commit if it is possible
func (o *Object) ParseCommit() (*Commit, error) {
	if o.Type != TypeCommit {
		return nil, errors.New("Object is not commit")
	}
	return DeserializeCommit(o.Data), nil
}
