package documentstore

import (
	"fmt"
)

type ICollection interface {
	Put(doc Document) error
	Get(key string) (*Document, error)
	Delete(key string) bool
	List() []Document
}
type Collection struct {
	Documents map[string]*Document
	name      string
	cfg       CollectionConfig
}

type CollectionConfig struct {
	PrimaryKey string
}

// Потрібно перевірити що документ містить поле `{cfg.PrimaryKey}` типу `string`
func (c *Collection) Put(doc Document) error {
	if field, ok := doc.Fields[c.cfg.PrimaryKey]; ok {
		if field.Type == DocumentFieldTypeString {
			c.Documents[field.Value.(string)] = &doc
		} else {
			return fmt.Errorf("%w:%s", ErrUnsupportedDocumentField, "primary key must be a string")
		}
	} else {
		return fmt.Errorf("%w", ErrEmptyPrimaryKey)
	}
	return nil
}

func (c *Collection) Get(key string) (*Document, error) {
	if col, ok := c.Documents[key]; ok {
		return col, nil
	}
	return nil, fmt.Errorf("%w: %s", ErrDocumentNotFound, key)

}

func (c *Collection) Delete(key string) bool {
	if _, ok := c.Documents[key]; ok {
		delete(c.Documents, key)
		return true
	}
	return false
}

func (c *Collection) List() []Document {
	list := make([]Document, 0, len(c.Documents))
	for _, doc := range c.Documents {
		list = append(list, *doc)
	}
	return list
}
