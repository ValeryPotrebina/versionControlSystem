package storage

import (
	"bytes"
	"fmt"
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
	Branch string
	Refs   map[string][]byte
	Path   string //Path to directory
}

type CommitData struct {
	Hash   []byte
	Commit *object.Commit
}

// Initialize database in path, if root hash not specified, build init file system and add data to database.
func InitStorage(path string) (*Storage, error) {
	if _, err := os.Open(path); os.IsNotExist(err) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	// Disable badger logs
	opts := badger.DefaultOptions(path + "/.vcs").WithLogger(nil)
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}

	storage := &Storage{
		db,
		"",
		make(map[string][]byte, 0),
		path,
	}

	branch, err := storage.GetData([]byte(BRANCH_KEY))
	if err == badger.ErrKeyNotFound {
		fmt.Println("BRANCH not found. Initializing BRANCH...")
		//create init commit
		tree := object.Tree{
			Children: []object.Child{},
		}
		treeObj := tree.CreateObject()
		treeHash, treeData := treeObj.GetData()
		err := storage.SetData(treeHash, treeData)
		if err != nil {
			return nil, err
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
			return nil, err
		}
		// fmt.Print("COMMIT-ORIGIN: ", commit.Origin)
		//initialize branches
		storage.Branch = MASTER_BRANCH
		storage.Refs[MASTER_BRANCH] = commitHash

		err = storage.SetData([]byte(BRANCH_KEY), []byte(storage.Branch))
		if err != nil {
			return nil, err
		}

		err = storage.SetData([]byte(REFS_KEY), SerializeRefs(storage.Refs))
		if err != nil {
			return nil, err
		}
		return storage, nil
	}
	refs, err := storage.GetData([]byte(REFS_KEY))
	if err != nil {
		return nil, err
	}

	fmt.Printf("BRANCH found: %s\n", branch)
	storage.Branch = string(branch)
	storage.Refs = DeserializeRefs(refs)

	return storage, nil
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

func (s *Storage) CreateCommit(author string, description string) error {
	fs := InitFileSystem(s.Path)
	for _, obj := range fs.TreeMap {
		hash, data := obj.GetData()
		s.SetData(hash, data)
	}
	commit := object.Commit{
		Origin:      s.Refs[s.Branch],
		Tree:        fs.ROOT_HASH,
		Author:      []byte(author),
		Description: []byte(description),
		Time:        time.Now().Unix(),
	}
	commitObj := commit.CreateObject()
	commitHash, commitData := commitObj.GetData()
	err := s.SetData(commitHash, commitData)
	if err != nil {
		return err
	}
	s.Refs[s.Branch] = commitHash
	err = s.SetData([]byte(REFS_KEY), SerializeRefs(s.Refs))
	return err
}

// TODO поправить
func (s *Storage) GetCommits(branch string, count uint64) ([]*CommitData, error) {
	commits := make([]*CommitData, 0)
	if s.Refs[branch] == nil {
		return commits, fmt.Errorf("branch \"%s\" does not exist", branch)
	}
	//hash current commit in fol branch
	commitHash := s.Refs[branch]
	for i := uint64(0); !bytes.Equal(commitHash, []byte{}) || (count > 0 && i >= count); i++ {
		commitData, err := s.GetCommit(commitHash)
		if err != nil {
			return commits, err
		}
		commits = append(commits, commitData)
		// commithash assign prev commit hash
		commitHash = commitData.Commit.Origin
	}

	return commits, nil
}

func (s *Storage) GetCommit(hash []byte) (*CommitData, error) {
	commitObj, err := s.GetObject(hash)
	if err != nil {
		return nil, err
	}
	commit, err := commitObj.ParseCommit()
	commitData := CommitData{
		Hash:   hash,
		Commit: commit,
	}
	return &commitData, err
}

// find diffs between file system stored in database and real file system
func (s *Storage) Diffs() ([]*object.FileChange, error) {
	fileChange := make([]*object.FileChange, 0)
	fs := InitFileSystem(s.Path)
	cmp := object.Comparator{
		GetFunction1: s.GetObject,
		GetFunction2: fs.GetObject,
	}
	commitHash := s.Refs[s.Branch]
	commit, err := s.GetCommit(commitHash)
	if err != nil {
		return fileChange, err
	}

	fileChange, err = cmp.CompareTrees(commit.Commit.Tree, fs.ROOT_HASH)

	return fileChange, err
}
func (s *Storage) DiffsWithCommit(hash []byte) ([]*object.FileChange, error) {
	fileChange := make([]*object.FileChange, 0)
	fs := InitFileSystem(s.Path)
	cmp := object.Comparator{
		GetFunction1: s.GetObject,
		GetFunction2: fs.GetObject,
	}
	commitObj, err := s.GetObject(hash)
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
func (s *Storage) DiffsBetweenCommits(hash1 []byte, hash2 []byte) ([]*object.FileChange, error) {
	cmp := object.Comparator{
		GetFunction1: s.GetObject,
		GetFunction2: s.GetObject,
	}

	fileChange, err := cmp.CompareCommits(hash1, hash2)

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

func (s *Storage) GetBranches() []string {
	branches := make([]string, 0, len(s.Refs))
	for k := range s.Refs {
		branches = append(branches, k)
	}
	return branches
}

func (s *Storage) ChangeBranch(branch string) error {
	if s.Refs[branch] == nil {
		return fmt.Errorf("branch \"%s\" does not exist", branch)
	}
	s.Branch = branch
	return nil
}

func (s *Storage) CreateBranch(branch string) error {
	if s.Refs[branch] != nil {
		return fmt.Errorf("branch \"%s\" already exists", branch)
	}
	s.Refs[branch] = s.Refs[s.Branch]
	return nil
}
