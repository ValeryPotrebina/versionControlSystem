package storage

import (
	// "fmt"
	"fmt"
	"mymodule/internal/object"
	"os"
)

// Current file system state
type FileSystem struct {
	path      string                    //Path to vcs
	ROOT_HASH []byte                    //Hash of root object
	TreeMap   map[string]*object.Object //Map of object for current file system state
}

// Get stored in map object for key
func (fs *FileSystem) GetObject(key []byte) (*object.Object, error) {
	return fs.TreeMap[fmt.Sprintf("%x", key)], nil
}

func (fs *FileSystem) SetObject(key []byte, data *object.Object) {
	fs.TreeMap[fmt.Sprintf("%x", key)] = data
}

func InitFileSystem(path string) (*FileSystem, error) {
	fs := &FileSystem{
		path,
		[]byte{},
		make(map[string]*object.Object),
	}
	rootTree, err := fs.CreateTree(path)
	if err != nil {
		return nil, err
	}
	hash, err := rootTree.GetHash()
	if err != nil {
		return nil, err
	}
	fs.ROOT_HASH = hash
	// fmt.Println("ROOT_HASH: ", fs.ROOT_HASH)
	return fs, nil
}

// Creating tree for current file system state
func (fs *FileSystem) CreateTree(path string) (*object.Object, error) {
	stat, err := os.Stat(path)
	if err != nil {

		return nil, err
	}

	//if path is file create Blob object
	if !stat.IsDir() {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		blob := object.Blob{
			Data: data,
		}
		obj := blob.CreateObject()

		hash, err := obj.GetHash()
		if err != nil {
			return nil, err
		}
		fs.SetObject(hash, obj)
		return obj, nil
	}

	//else collect children objects
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	children := make([]object.Child, 0)
	for _, e := range entries {
		if e.Name() == ".vcs" {
			continue
		}
		obj, err := fs.CreateTree(path + "/" + e.Name())
		if err != nil {
			return nil, err
		}
		hash, err := obj.GetHash()
		if err != nil {
			return nil, err
		}
		children = append(children, object.Child{
			Type: obj.Type,
			Name: []byte(e.Name()),
			Hash: hash,
		})
	}
	//and create tree object with that children
	tree := object.Tree{
		Children: children,
	}
	obj, err := tree.CreateObject()
	if err != nil {
		return nil, err
	}
	hash, err := obj.GetHash()
	if err != nil {
		return nil, err
	}
	fs.SetObject(hash, obj)
	return obj, nil
}
