package documentstore

import (
	"fmt"
	"sync"

	"github.com/google/btree"
)

type QueryParams struct {
	Desc     bool    // Визначає в якому порядку повертати дані
	MinValue *string // Визначає мінімальне значення поля для фільтрації
	MaxValue *string // Визначає максимальне значення поля для фільтрації
}

type Collector interface {
	Put(doc Document) error
	Get(key string) (*Document, error)
	Delete(key string) bool
	List() []Document
	Query(fieldName string, params QueryParams) ([]Document, error)
}
type Collection struct {
	Documents map[string]*Document
	name      string
	Cfg       CollectionConfig
	Indexes   map[string]*Index
	mu        sync.RWMutex
}

func NewCollection(name string, Cfg *CollectionConfig) *Collection {
	return &Collection{
		Documents: make(map[string]*Document),
		name:      name,
		Cfg:       *Cfg,
		Indexes:   make(map[string]*Index),
		mu:        sync.RWMutex{},
	}
}

func (c *Collection) CreateIndex(fieldName string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	fmt.Println("Create index for field: ", fieldName)
	if _, ok := c.Indexes[fieldName]; ok {
		return fmt.Errorf("index for field '%s' already exists", fieldName)
	}

	index := NewIndex(fieldName)

	for _, doc := range c.Documents {
		if err := index.Put(doc, c.Cfg.PrimaryKey); err != nil {
			return err
		}
	}
	c.Indexes[fieldName] = index

	return nil
}

func (c *Collection) DeleteIndex(fieldName string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.Indexes[fieldName]; !ok {
		return fmt.Errorf("index for field '%s' doesn't exist", fieldName)
	}
	delete(c.Indexes, fieldName)
	return nil
}

type CollectionConfig struct {
	PrimaryKey string
}

func queryHelper(params QueryParams, res *[]Document) func(item btree.Item) bool {
	return func(item btree.Item) bool {
		i := item.(Item)

		if params.MinValue != nil && i.Value < *params.MinValue {
			return true
		}
		if params.MaxValue != nil && i.Value > *params.MaxValue {
			return false
		}
		*res = append(*res, *i.Doc)
		return true
	}
}

// Якщо для даного поля не існує індексу - повертаємо помилку
func (c *Collection) Query(fieldName string, params QueryParams) ([]Document, error) {
	c.mu.RLock()
	index, ok := c.Indexes[fieldName]
	c.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("index for field '%s' doesn't exist", fieldName)
	}
	list := make([]Document, 0)
	cb := queryHelper(params, &list)

	index.mu.RLock()
	defer index.mu.RUnlock()
	//1. MinValue == nil && MaxValue == nil + ASC & DESC
	if params.MinValue == nil && params.MaxValue == nil {
		if params.Desc {
			index.Tree.Descend(cb)
		} else {
			index.Tree.Ascend(cb)
		}
		return list, nil
	}

	//2. MinValue != nil && MaxValue != nil + ASC & DESC
	if params.MinValue != nil && params.MaxValue != nil {
		minValue := Item{Value: *params.MinValue, Id: ""}
		maxValue := Item{Value: *params.MaxValue, Id: "\xFF"}
		if params.Desc {
			index.Tree.DescendRange(maxValue, minValue, cb)
		} else {
			index.Tree.AscendRange(minValue, maxValue, cb)
		}
		return list, nil
	}

	//3. MaxValue != nil && MinValue == nil + ASC & DESC
	if params.MinValue != nil {
		minValue := Item{Value: *params.MinValue, Id: ""}

		if params.Desc {
			index.Tree.DescendGreaterThan(minValue, cb)
		} else {
			index.Tree.AscendGreaterOrEqual(minValue, cb)
		}
		return list, nil
	}

	//4. MaxValue == nil && MinValue != nil + ASC & DESC
	maxValue := Item{Value: *params.MaxValue, Id: "\xFF"}

	if params.Desc {
		index.Tree.DescendLessOrEqual(maxValue, cb)
	} else {
		index.Tree.AscendLessThan(maxValue, cb)
	}

	return list, nil
}

// Потрібно перевірити що документ містить поле `{Cfg.PrimaryKey}` типу `string`
func (c *Collection) Put(doc Document) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	field, ok := doc.Fields[c.Cfg.PrimaryKey]
	if ok {
		if field.Type == DocumentFieldTypeString {

			if d, exists := c.Documents[field.Value.(string)]; exists {
				for _, index := range c.Indexes {
					err := index.Delete(d, c.Cfg.PrimaryKey)
					if err != nil {
						return err
					}
				}
			}

			c.Documents[field.Value.(string)] = &doc
			for _, index := range c.Indexes {
				err := index.Put(&doc, c.Cfg.PrimaryKey)
				if err != nil {
					return err
				}
			}
			return nil
		}
		return fmt.Errorf("%w:%s", ErrUnsupportedDocumentField, "primary key must be a string")
	}
	return fmt.Errorf("%w", ErrEmptyPrimaryKey)
}

func (c *Collection) Get(key string) (*Document, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if col, ok := c.Documents[key]; ok {
		return col, nil
	}
	return nil, fmt.Errorf("%w: %s", ErrDocumentNotFound, key)

}

func (c *Collection) Delete(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if doc, ok := c.Documents[key]; ok {
		for _, index := range c.Indexes {
			err := index.Delete(doc, c.Cfg.PrimaryKey)
			if err != nil {
				return false
			}
		}
		delete(c.Documents, key)
		return true
	}
	return false
}

func (c *Collection) List() []Document {
	c.mu.RLock()
	defer c.mu.RUnlock()
	list := make([]Document, 0, len(c.Documents))
	for _, doc := range c.Documents {
		list = append(list, *doc)
	}
	return list
}
