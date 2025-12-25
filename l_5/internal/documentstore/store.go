package documentstore

import (
	"errors"
	"fmt"
)

var ErrCollectionAlreadyExists = errors.New("collection already exists")
var ErrCollectionNotFound = errors.New("collection not found")
var ErrDocumentNotFound = errors.New("document not found")
var ErrUnsupportedDocumentField = errors.New("unsupported document field")
var ErrEmptyCollectionName = errors.New("empty collection name")
var ErrEmptyPrimaryKey = errors.New("empty primary key")

type Store struct {
	Collections map[string]*Collection
}

func NewStore() *Store {
	return &Store{
		Collections: make(map[string]*Collection),
	}
}

func (s *Store) CreateCollection(name string, cfg *CollectionConfig) (*Collection, error) {
	// Створюємо нову колекцію і повертаємо `true` якщо колекція була створена
	// Якщо ж колекція вже створеня то повертаємо `false` та nil
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
