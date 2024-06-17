package storage

import (
	// "fmt"
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
func (fs *FileSystem) GetData(key []byte) (*object.Object, error) {
	return fs.TreeMap[string(key)], nil
}

func (fs *FileSystem) SetData(key []byte, data *object.Object) {
	fs.TreeMap[string(key)] = data
}

func InitFileSystem(path string) *FileSystem {
	fs := FileSystem{
		path,
		// Почему не вычисляем рутовский хэш тут
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
	// Скорее тогда создаем дерево объетов
	// Почему везде звездочки...
	// Какой путь передаем в функцию 
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

		fs.SetData(obj.GetHash(), obj)
		return obj
		// После того как ретерним, дальше будет разюор дерева?
		// Как происходит разбор 
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

	fs.SetData(obj.GetHash(), obj)
	return obj
}

