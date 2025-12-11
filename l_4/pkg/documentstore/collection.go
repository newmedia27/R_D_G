package documentstore

import "fmt"

type Collection struct {
	Documents map[string]*Document
	name      string
	cfg       CollectionConfig
}

type CollectionConfig struct {
	PrimaryKey string
}

func (c *Collection) Put(doc Document) {
	// Потрібно перевірити що документ містить поле `{cfg.PrimaryKey}` типу `string`
	if field, ok := doc.Fields[c.cfg.PrimaryKey]; ok {
		if field.Type == DocumentFieldTypeString {
			c.Documents[field.Value.(string)] = &doc
		} else {
			fmt.Println("Primary key must be string")
		}
	} else {
		fmt.Println("Document doesn't have primary key")
	}
}

func (c *Collection) Get(key string) (*Document, bool) {

	if _, ok := c.Documents[key]; ok {
		return c.Documents[key], true
	}
	return nil, false

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
