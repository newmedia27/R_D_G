package documentstore

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var ErrCollectionAlreadyExists = errors.New("collection already exists")
var ErrCollectionNotFound = errors.New("collection not found")
var ErrDocumentNotFound = errors.New("document not found")
var ErrUnsupportedDocumentField = errors.New("unsupported document field")
var ErrEmptyCollectionName = errors.New("empty collection name")
var ErrEmptyPrimaryKey = errors.New("empty primary key")

type Dump interface {
	Dump() ([]byte, error)
	DumpToFIle(filename string) error
}
type Store struct {
	Collections map[string]*Collection
}

func NewStore() *Store {
	return &Store{
		Collections: make(map[string]*Collection),
	}
}

func NewStoreFromDump(dump []byte) (*Store, error) {
	serialize := make(map[string]*Collection)

	err := json.Unmarshal(dump, &serialize)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", err, "can't deserialize dump")
	}
	return &Store{
		Collections: serialize,
	}, nil
}

func NewStoreFromFile(filename string) (*Store, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", err, "can't open file")
	}
	defer func() {
		err = file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", err, "can't read file")
	}
	return NewStoreFromDump(data)
}

// Створюємо нову колекцію і повертаємо `true` якщо колекція була створена
// Якщо ж колекція вже створеня  то повертаємо `false` та nil
func (s *Store) CreateCollection(name string, cfg *CollectionConfig) (*Collection, error) {
	if name == "" {
		return nil, fmt.Errorf("%w: %s", ErrEmptyCollectionName, "collection name can't be empty")
	}
	if _, ok := s.Collections[name]; ok {
		return nil, fmt.Errorf("%w: %s", ErrCollectionAlreadyExists, name)
	}

	col := &Collection{
		Documents: make(map[string]*Document),
		name:      name,
		cfg:       *cfg,
	}
	s.Collections[name] = col

	return col, nil
}

func (s *Store) GetCollection(name string) (*Collection, error) {
	if col, ok := s.Collections[name]; ok {
		return col, nil
	}
	return nil, fmt.Errorf("%w: %s", ErrCollectionNotFound, name)
}

func (s *Store) DeleteCollection(name string) bool {
	if _, ok := s.Collections[name]; ok {
		delete(s.Collections, name)
		return true
	}
	return false
}

func (s *Store) Dump() ([]byte, error) {
	dump, err := json.Marshal(s.Collections)
	if err != nil {
		return nil, err
	}
	return dump, nil
}

func checkDir(path string) (string, error) {
	dir := filepath.Dir(path)
	if dir != "." {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return "", fmt.Errorf("can't create directory %s: %w", dir, err)
		}
	}
	return path, nil
}

func (s *Store) DumpToFIle(path string) error {
	dump, err := s.Dump()
	if err != nil {
		return err
	}
	filename, err := checkDir(path)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer func() {
		err = file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()
	_, err = file.Write(dump)
	if err != nil {
		return err
	}
	err = file.Sync()
	if err != nil {
		return err
	}
	return nil
}
