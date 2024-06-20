package object

import (
	"bytes"
	"encoding/gob"
	"errors"
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
func (o *Object) GetHash() ([]byte, error) {
	data, err := o.Serialize()
	if err != nil {
		return nil, err
	}
	return CalculateHash(data), nil
}

// Get pair of hash and zipped data of object
func (o *Object) GetData() (hash []byte, data []byte, err error) {
	var b []byte
	b, err = o.Serialize()
	if err != nil {
		return
	}
	hash = CalculateHash(b)
	data, err = Zip(b)
	return
}

// Serialize object
func (o *Object) Serialize() ([]byte, error) {
	var b bytes.Buffer
	encoder := gob.NewEncoder(&b)
	err := encoder.Encode(o)
	return b.Bytes(), err
}

// Deserialize data into object
func DeserializeObject(data []byte) (*Object, error) {
	var obj Object
	unzipData, err := Unzip(data)
	if err != nil {
		return nil, err
	}
	decoder := gob.NewDecoder(bytes.NewReader(unzipData))
	err = decoder.Decode(&obj)
	if err != nil {
		return nil, err
	}
	return &obj, nil
}

// Unpack object into tree if it is possible
func (o *Object) ParseTree() (*Tree, error) {
	if o.Type != TypeTree {
		return nil, errors.New("Object is not tree")
	}
	return DeserializeTree(o.Data)
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
	return DeserializeCommit(o.Data)

}
