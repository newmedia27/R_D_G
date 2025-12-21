package documentstore

import (
	"errors"

	"github.com/newmedia27/R_D_G/l_3/pkg/printobject"
)

type DocumentFieldType string

const (
	DocumentFieldTypeString DocumentFieldType = "string"
	DocumentFieldTypeNumber DocumentFieldType = "number"
	DocumentFieldTypeBool   DocumentFieldType = "bool"
	DocumentFieldTypeArray  DocumentFieldType = "array"
	DocumentFieldTypeObject DocumentFieldType = "object"
)

type DocumentField struct {
	Type  DocumentFieldType
	Value interface{}
}

type Document struct {
	Fields map[string]DocumentField
}

var documents = map[string]*Document{}

func Put(doc *Document) error {
	// 1. Перевірити що документ містить в мапі поле `key` типу `string`
	// 2. Додати Document до локальної мапи з документами

	key, isIsset := doc.Fields["key"]
	if !isIsset || key.Type != DocumentFieldTypeString {
		return errors.New("document must contain 'key' field and field must be of string type")
	}

	documents[key.Value.(string)] = doc

	printobject.PrintObject("Document added:", documents) // TODO: тут питання, моя основна мова це - JS, і там ну дуже класний log під капотом,  а як тут бути, особливо якщо скажена вкладеність??? типу slog/logrus/zerolog etc.??
	return nil
}

func Get(key string) (*Document, bool) {
	// Потрібно повернути документ по ключу
	// Якщо документ знайдено, повертаємо `true` та поінтер на документ
	// Інакше повертаємо `false` та `nil`
	doc, isIsset := documents[key]
	return doc, isIsset
}

func Delete(key string) bool {
	// Видаляємо документа по ключу.
	// Повертаємо `true` якщо ми знайшли і видалили документі
	// Повертаємо `false` якщо документ не знайдено
	_, isIsset := documents[key]
	if isIsset {
		delete(documents, key)
		return true
	}

	return false
}

func List() (docs []*Document) {
	// Повертаємо список усіх документів
	for _, doc := range documents {
		docs = append(docs, doc)
	}
	return
}
