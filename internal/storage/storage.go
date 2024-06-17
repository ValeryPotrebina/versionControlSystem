package storage

import (
	"bytes"
	"log"
	"mymodule/internal/object"
	"os"
	"time"

	"github.com/dgraph-io/badger"
)

const BRANCH_KEY = "BRANCH"
const REFS_KEY = "REFS"

const INITIAL_COMMIT = "Initial commit"
const MASTER_BRANCH = "master"

// Storage of version control system
type Storage struct {
	DB     *badger.DB //Database for object storing
	branch string
	refs   map[string][]byte
	Path   string //Path to directory
}

// Initialize database in path, if root hash not specified, build init file system and add data to database.
func InitStorage(path string) Storage {
	if _, err := os.Open(path); os.IsNotExist(err) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Panic(err)
		}
	}

	opts := badger.DefaultOptions(path + "/.vcs")
	db, err := badger.Open(opts)
	if err != nil {
		log.Panic(err)
	}

	storage := Storage{
		db,
		"",
		make(map[string][]byte, 0),
		path,
	}

	branch, err := storage.GetData([]byte(BRANCH_KEY))
	if err == badger.ErrKeyNotFound {
		log.Println("BRANCH not found. Initializing BRANCH...")
		//create init commit
		tree := object.Tree{
			Children: []object.Child{},
		}
		treeObj := tree.CreateObject()
		treeHash, treeData := treeObj.GetData()
		err := storage.SetData(treeHash, treeData)
		if err != nil {
			log.Panic(err)
		}
		commit := object.Commit{
			Origin:      []byte{},
			Tree:        treeHash,
			Author:      []byte{},
			Time:        time.Now().Unix(),
			Description: []byte(INITIAL_COMMIT),
		}
		commitObj := commit.CreateObject()
		commitHash, commitData := commitObj.GetData()
		err = storage.SetData(commitHash, commitData)
		if err != nil {
			log.Panic(err)
		}
		// fmt.Print("COMMIT-ORIGIN: ", commit.Origin)
		//initialize branches
		storage.branch = MASTER_BRANCH
		storage.refs[MASTER_BRANCH] = commitHash

		err = storage.SetData([]byte(BRANCH_KEY), []byte(storage.branch))
		if err != nil {
			log.Panic(err)
		}
		//Не понятна шняга с ветками и рефами
		err = storage.SetData([]byte(REFS_KEY), SerializeRefs(storage.refs))
		if err != nil {
			log.Panic(err)
		}
		return storage
	}
	refs, err := storage.GetData([]byte(REFS_KEY))
	if err != nil {
		log.Panic(err)
	}

	log.Printf("BRANCH found: %x", branch)
	storage.branch = string(branch)
	storage.refs = DeserializeRefs(refs)

	return storage
}

// Get data from database for this key
// Что является ключом то....то бранч, то что-то еще судя по коду
func (s *Storage) GetData(key []byte) ([]byte, error) {
	var hash []byte

	err := s.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			hash = val
			return nil
		})
		return err
	})

	return hash, err
}

// Set data for key in database
func (s *Storage) SetData(key []byte, data []byte) error {
	err := s.DB.Update(func(txn *badger.Txn) error {
		err := txn.Set(key, data)
		return err
	})
	return err
}

// Close database
func (s *Storage) CloseStorage() {
	s.DB.Close()
}

func (s *Storage) CreateCommit(author string, description string) {
	fs := InitFileSystem(s.Path)
	for _, obj := range fs.TreeMap {
		hash, data := obj.GetData()
		s.SetData(hash, data)
	}
	commit := object.Commit{
		Origin:      s.refs[s.branch],
		Tree:        fs.ROOT_HASH,
		Author:      []byte(author),
		Description: []byte(description),
		Time:        time.Now().Unix(),
	}
	commitObj := commit.CreateObject()
	commitHash, commitData := commitObj.GetData()
	err := s.SetData(commitHash, commitData)
	if err != nil {
		log.Panic(err)
	}
	s.refs[s.branch] = commitHash
	err = s.SetData([]byte(REFS_KEY), SerializeRefs(s.refs))
	if err != nil {
		log.Panic(err)
	}
}

// TODO поправить
func (s *Storage) GetCommits(count int) map[string]*object.Commit {
	commits := make(map[string]*object.Commit, 0)
	commitHash := s.refs[s.branch]
	for i := 0; !bytes.Equal(commitHash, []byte{}); i++ {
		if count > 0 && i >= count {
			break
		}
		objData, err := s.GetData(commitHash)
		if err != nil {
			log.Panic(err)
		}
		commitObj := object.DeserializeObject(objData)
		commit, err := commitObj.ParseCommit()
		if err != nil {
			log.Panic(err)
		}
		commits[string(commitHash)] = commit
		commitHash = commit.Origin
	}

	return commits
}

func (s *Storage) GetCommit(hash []byte) (*object.Commit, error) {
	commitObj, err := s.GetObject(hash)
	if err != nil {
		return nil, err
	}
	commit, err := commitObj.ParseCommit()
	return commit, err
}

// find diffs between file system stored in database and real file system
func (s *Storage) FindDiffs() ([]*object.FileChange, error) {
	fileChange := make([]*object.FileChange, 0)
	fs := InitFileSystem(s.Path)
	cmp := object.Comparator{
		GetFunction1: s.GetObject,
		GetFunction2: fs.GetData,
	}
	commitHash := s.refs[s.branch]
	commitObj, err := s.GetObject(commitHash)
	if err != nil {
		return fileChange, err
	}
	commit, err := commitObj.ParseCommit()
	if err != nil {
		return fileChange, err
	}

	fileChange, err = cmp.CompareTrees(commit.Tree, fs.ROOT_HASH)

	return fileChange, err
}

func (s *Storage) GetObject(key []byte) (*object.Object, error) {
	objData, err := s.GetData(key)
	if err != nil {
		return nil, err
	}
	obj := object.DeserializeObject(objData)
	return obj, nil
}

func (s *Storage) SetObject(obj *object.Object) error {
	objKey, objData := obj.GetData()
	err := s.SetData(objKey, objData)
	return err
}

func (s *Storage) GetBranches() ([]string, string) {
	branches := make([]string, 0, len(s.refs))
	for k := range s.refs {
		branches = append(branches, k)
	}
	return branches, s.branch
}
