package object

import (
	"bytes"

	"github.com/sergi/go-diff/diffmatchpatch"
)

const (
	ActionInsert = iota
	ActionDelete
)

type Comparator struct {
	GetFunction1 func([]byte) (*Object, error)
	GetFunction2 func([]byte) (*Object, error)
}

type FileChange struct {
	FileName []byte
	Changes  []diffmatchpatch.Diff
}

func (cmp *Comparator) CompareTrees(hash1 []byte, hash2 []byte) ([]*FileChange, error) {
	fileChanges := make([]*FileChange, 0)

	if bytes.Equal(hash1, hash2) {
		return fileChanges, nil
	}

	var obj1, obj2 *Object
	var t1, t2 *Tree
	var err error
	if hash1 != nil {
		obj1, err = cmp.GetFunction1(hash1)
		if err != nil {
			return fileChanges, err
		}
		t1, err = obj1.ParseTree()
		if err != nil {
			return fileChanges, err
		}
	}

	if hash2 != nil {
		obj2, err = cmp.GetFunction2(hash2)
		if err != nil {
			return fileChanges, err
		}
		t2, err = obj2.ParseTree()
		if err != nil {
			return fileChanges, err
		}
	}

	if hash1 == nil || hash2 == nil {
		if hash1 == nil {
			for _, c := range t2.Children {
				switch c.Type {
				case TypeBlob:
					diffs, err := cmp.CompareBlobs(nil, c.Hash)
					if err != nil {
						return fileChanges, err
					}
					if len(diffs) > 0 {
						fileChanges = append(fileChanges, &FileChange{
							FileName: c.Name,
							Changes:  diffs,
						})
					}

				case TypeTree:
					changes, err := cmp.CompareTrees(nil, c.Hash)
					if err != nil {
						return fileChanges, err
					}
					for _, change := range changes {
						change.FileName = bytes.Join([][]byte{c.Name, []byte("/"), change.FileName}, []byte(""))
						fileChanges = append(fileChanges, change)
					}
				}

			}
		}

		if hash2 == nil {
			for _, c := range t1.Children {
				switch c.Type {
				case TypeBlob:
					diffs, err := cmp.CompareBlobs(c.Hash, nil)
					if err != nil {
						return fileChanges, err
					}
					if len(diffs) > 0 {
						fileChanges = append(fileChanges, &FileChange{
							FileName: c.Name,
							Changes:  diffs,
						})
					}
				case TypeTree:
					changes, err := cmp.CompareTrees(c.Hash, nil)
					if err != nil {
						return fileChanges, err
					}
					for _, change := range changes {
						change.FileName = bytes.Join([][]byte{c.Name, []byte("/"), change.FileName}, []byte(""))
						fileChanges = append(fileChanges, change)
					}
				}

			}
		}

		return fileChanges, nil
	}

	for _, c1 := range t1.Children {
		found := false
		for _, c2 := range t2.Children {
			if bytes.Equal(c1.Name, c2.Name) {
				found = true
				if c1.Type == c2.Type {
					switch c1.Type {
					case TypeBlob:
						diffs, err := cmp.CompareBlobs(c1.Hash, c2.Hash)
						if err != nil {
							return fileChanges, err
						}
						if len(diffs) > 0 {
							fileChanges = append(fileChanges, &FileChange{
								FileName: c1.Name,
								Changes:  diffs,
							})
						}

					case TypeTree:
						changes, err := cmp.CompareTrees(c1.Hash, c2.Hash)
						if err != nil {
							return fileChanges, err
						}
						for _, change := range changes {
							change.FileName = bytes.Join([][]byte{c1.Name, []byte("/"), change.FileName}, []byte(""))
							fileChanges = append(fileChanges, change)
						}

					}

				} else {
					switch c1.Type {
					case TypeBlob:
						diffs, err := cmp.CompareBlobs(c1.Hash, nil)
						if err != nil {
							return fileChanges, err
						}
						if len(diffs) > 0 {
							fileChanges = append(fileChanges, &FileChange{
								FileName: c1.Name,
								Changes:  diffs,
							})
						}
					case TypeTree:
						changes, err := cmp.CompareTrees(c1.Hash, nil)
						if err != nil {
							return fileChanges, err
						}
						for _, change := range changes {
							change.FileName = bytes.Join([][]byte{c1.Name, []byte("/"), change.FileName}, []byte(""))
							fileChanges = append(fileChanges, change)
						}
					}

					switch c2.Type {
					case TypeBlob:
						diffs, err := cmp.CompareBlobs(nil, c2.Hash)
						if err != nil {
							return fileChanges, err
						}
						if len(diffs) > 0 {
							fileChanges = append(fileChanges, &FileChange{
								FileName: c2.Name,
								Changes:  diffs,
							})
						}
					case TypeTree:
						changes, err := cmp.CompareTrees(nil, c2.Hash)
						if err != nil {
							return fileChanges, err
						}
						for _, change := range changes {
							change.FileName = bytes.Join([][]byte{c2.Name, []byte("/"), change.FileName}, []byte(""))
							fileChanges = append(fileChanges, change)
						}
					}
				}
				break
			}
		}
		// Если раньше был ребенок, а теперь его нет (удаление файла, папки)
		if !found {
			switch c1.Type {
			case TypeBlob:
				diffs, err := cmp.CompareBlobs(c1.Hash, nil)
				if err != nil {
					return fileChanges, err
				}
				if len(diffs) > 0 {
					fileChanges = append(fileChanges, &FileChange{
						FileName: c1.Name,
						Changes:  diffs,
					})
				}
			case TypeTree:
				changes, err := cmp.CompareTrees(c1.Hash, nil)
				if err != nil {
					return fileChanges, err
				}
				for _, change := range changes {
					change.FileName = bytes.Join([][]byte{c1.Name, []byte("/"), change.FileName}, []byte(""))
					fileChanges = append(fileChanges, change)
				}
			}
		}
	}

	for _, c2 := range t2.Children {
		found := false
		for _, c1 := range t1.Children {
			if bytes.Equal(c2.Name, c1.Name) {
				found = true
				break
			}
		}
		if !found {
			switch c2.Type {
			case TypeBlob:
				diffs, err := cmp.CompareBlobs(nil, c2.Hash)
				if err != nil {
					return fileChanges, err
				}
				if len(diffs) > 0 {
					fileChanges = append(fileChanges, &FileChange{
						FileName: c2.Name,
						Changes:  diffs,
					})
				}
			case TypeTree:
				changes, err := cmp.CompareTrees(nil, c2.Hash)
				if err != nil {
					return fileChanges, err
				}
				for _, change := range changes {
					change.FileName = bytes.Join([][]byte{c2.Name, []byte("/"), change.FileName}, []byte(""))
					fileChanges = append(fileChanges, change)
				}
			}
		}
	}
	return fileChanges, nil
}

func (cmp *Comparator) CompareBlobs(hash1 []byte, hash2 []byte) ([]diffmatchpatch.Diff, error) {
	diffs := make([]diffmatchpatch.Diff, 0)
	if bytes.Equal(hash1, hash2) {
		return diffs, nil
	}
	var obj1, obj2 *Object
	var b1, b2 *Blob
	var err error

	if hash1 != nil {
		obj1, err = cmp.GetFunction1(hash1)
		if err != nil {
			return diffs, err
		}
		b1, err = obj1.ParseBlob()
		if err != nil {
			return diffs, err
		}
	}
	if hash2 != nil {
		obj2, err = cmp.GetFunction2(hash2)
		if err != nil {
			return diffs, err
		}
		b2, err = obj2.ParseBlob()
		if err != nil {
			return diffs, err
		}
	}

	dmp := diffmatchpatch.New()
	source1, source2 := "", ""
	if hash1 != nil {
		source1 = string(b1.Data)
	}
	if hash2 != nil {
		source2 = string(b2.Data)
	}
	lines1, lines2, lines := dmp.DiffLinesToChars(source1, source2)
	diffs = dmp.DiffMain(lines1, lines2, false)
	diffs = dmp.DiffCharsToLines(diffs, lines)

	return diffs, nil
}

func (cmp *Comparator) CompareCommits(hash1 []byte, hash2 []byte) ([]*FileChange, error) {
	fileChanges := make([]*FileChange, 0)
	if bytes.Equal(hash1, hash2) {
		return fileChanges, nil
	}
	obj1, err := cmp.GetFunction1(hash1)
	if err != nil {
		return fileChanges, err
	}
	obj2, err := cmp.GetFunction2(hash2)
	if err != nil {
		return fileChanges, err
	}
	c1, err := obj1.ParseCommit()
	if err != nil {
		return fileChanges, err
	}
	c2, err := obj2.ParseCommit()
	if err != nil {
		return fileChanges, err
	}
	fileChanges, err = cmp.CompareTrees(c1.Tree, c2.Tree)

	return fileChanges, err
}
