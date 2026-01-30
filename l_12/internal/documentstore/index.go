package documentstore

import (
	"errors"
	"sync"

	"github.com/google/btree"
)

type Indexer interface {
	Put(doc *Document, key string) error
	Delete(doc *Document, key string) error
}

type Item struct {
	Id    string
	Value string
	Doc   *Document
}

func (i Item) Less(b btree.Item) bool {
	other := b.(Item)
	if i.Value != other.Value {
		return i.Value < other.Value
	}
	return i.Id < other.Id
}

type Index struct {
	Name string       `json:"-"`
	Tree *btree.BTree `json:"-"`
	mu   sync.RWMutex
}

func NewIndex(name string) *Index {
	return &Index{
		Name: name,
		Tree: btree.New(10),
	}
}

func (i *Index) Put(doc *Document, key string) error {
	keyField, ok := doc.Fields[key]
	if !ok {
		return ErrEmptyPrimaryKey
	}

	id, ok := keyField.Value.(string)
	if !ok {
		return errors.New("primary key must be a string")
	}

	valueField, exists := doc.Fields[i.Name]
	// не помилка, просто пропускаємо за умовою задачі TODO:
	if !exists || valueField.Type != DocumentFieldTypeString {
		return nil
	}

	i.Tree.ReplaceOrInsert(Item{
		Value: valueField.Value.(string),
		Id:    id,
		Doc:   doc,
	})
	return nil
}

func (i *Index) Delete(doc *Document, key string) error {
	var id string

	keyField, ok := doc.Fields[key]
	if !ok {
		return nil
	}
	id, ok = keyField.Value.(string)
	if !ok {
		return nil
	}
	valueField, exists := doc.Fields[i.Name]
	if !exists || valueField.Type != DocumentFieldTypeString {
		return nil
	}

	i.Tree.Delete(Item{
		Value: valueField.Value.(string),
		Id:    id,
		Doc:   doc,
	})
	return nil
}
