package documentstore

import (
	"encoding/json"
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

func TypeValidate(v any) DocumentFieldType {
	switch v.(type) {
	case string:
		return DocumentFieldTypeString
	case int:
		return DocumentFieldTypeNumber
	case bool:
		return DocumentFieldTypeBool
	case []interface{}:
		return DocumentFieldTypeArray
	case map[string]interface{}:
		return DocumentFieldTypeObject
	default:
		return DocumentFieldTypeString
	}
}

func MarshalDocument(input any) (*Document, error) {

	j, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	var res map[string]interface{}
	res = make(map[string]interface{})
	err = json.Unmarshal(j, &res)
	if err != nil {
		return nil, err
	}
	doc := &Document{
		Fields: make(map[string]DocumentField),
	}
	for k, v := range res {
		doc.Fields[k] = DocumentField{
			Type:  TypeValidate(v),
			Value: v,
		}
	}
	return doc, err
}

func UnmarshalDocument(doc *Document, output any) error {
	m := make(map[string]interface{})
	for k, v := range doc.Fields {
		m[k] = v.Value
	}
	j, err := json.Marshal(m)
	if err != nil {
		return err
	}
	err = json.Unmarshal(j, output)
	return err
}
