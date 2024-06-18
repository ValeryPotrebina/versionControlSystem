package storage

import (
	// "fmt"
	"fmt"
	"log"
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

func InitFileSystem(path string) *FileSystem {
	fs := FileSystem{
		path,
		[]byte{},
		make(map[string]*object.Object),
	}
	rootTree := fs.CreateTree(path)
	fs.ROOT_HASH = rootTree.GetHash()
	// fmt.Println("ROOT_HASH: ", fs.ROOT_HASH)
	return &fs
}

// Creating tree for current file system state
func (fs *FileSystem) CreateTree(path string) *object.Object {
	stat, err := os.Stat(path)
	if err != nil {
		log.Panic(err)
	}

	//if path is file create Blob object
	if !stat.IsDir() {
		data, err := os.ReadFile(path)
		if err != nil {
			log.Panic(err)
		}
		blob := object.Blob{
			Data: data,
		}
		obj := blob.CreateObject()

		fs.SetObject(obj.GetHash(), obj)
		return obj
	}

	//else collect children objects
	entries, err := os.ReadDir(path)
	if err != nil {
		log.Panic(err)
	}
	children := make([]object.Child, 0)
	for _, e := range entries {
		if e.Name() == ".vcs" {
			continue
		}
		obj := fs.CreateTree(path + "/" + e.Name())
		children = append(children, object.Child{
			Type: obj.Type,
			Name: []byte(e.Name()),
			Hash: obj.GetHash(),
		})
	}
	//and create tree object with that children
	tree := object.Tree{
		Children: children,
	}
	obj := tree.CreateObject()

	fs.SetObject(obj.GetHash(), obj)
	return obj
}
