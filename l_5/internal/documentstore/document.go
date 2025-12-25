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

// Грався в reflect
//func TypeValidate(value reflect.Kind) (DocumentFieldType, error) {
//	switch value {
//	case reflect.String:
//		return DocumentFieldTypeString, nil
//	case reflect.Int:
//		return DocumentFieldTypeNumber, nil
//	case reflect.Bool:
//		return DocumentFieldTypeBool, nil
//	case reflect.Array, reflect.Slice:
//		return DocumentFieldTypeArray, nil
//	case reflect.Map:
//		return DocumentFieldTypeObject, nil
//	default:
//		return "", fmt.Errorf("type validate error: %w", ErrUnsupportedDocumentField)
//	}
//}

//func MarshalDocument(input any) (*Document, error) {
//	t := reflect.TypeOf(input)
//	v := reflect.ValueOf(input)
//	if t.Kind() != reflect.Struct && t.Kind() != reflect.Ptr {
//		return nil, fmt.Errorf("input must be a struct")
//	}
//	if t.Kind() == reflect.Ptr {
//		v = v.Elem()
//		t = t.Elem()
//	}
//	doc := &Document{
//		Fields: make(map[string]DocumentField),
//	}
//	for i := 0; i < t.NumField(); i++ {
//		field := t.Field(i)
//		value := v.Field(i)
//		fieldName := field.Tag.Get("json")
//		if fieldName == "" {
//			fieldName = field.Name
//		}
//		if type_, err := TypeValidate(value.Kind()); err != nil {
//			return nil, err
//		} else {
//			doc.Fields[fieldName] = DocumentField{
//				Type:  type_,
//				Value: value.Interface(),
//			}
//		}
//	}
//	return doc, nil
//}
//
//func UnmarshalDocument(doc *Document, output any) error {
//	outT := reflect.TypeOf(output)
//	docT := reflect.TypeOf(doc)
//	if outT.Kind() != reflect.Ptr {
//		fmt.Println("output must be a pointer")
//		return fmt.Errorf("output must be a pointer")
//	}
//	if docT.Kind() != reflect.Ptr {
//		fmt.Println("doc must be a pointer")
//		return fmt.Errorf("doc must be a pointer")
//	}
//	outT = outT.Elem()
//	outV := reflect.ValueOf(output).Elem()
//
//	for i := 0; i < outT.NumField(); i++ {
//		field := outT.Field(i)
//		value := outV.Field(i)
//		fieldName := field.Tag.Get("json")
//		if fieldName == "" {
//			fieldName = field.Name
//		}
//		docField, ok := doc.Fields[fieldName]
//		if !ok {
//			return fmt.Errorf("field %s not found", fieldName)
//		}
//		if type_, err := TypeValidate(value.Kind()); err != nil {
//			return err
//		} else if type_ != docField.Type {
//			return fmt.Errorf("field %s type mismatch", fieldName)
//		}
//
//		value.Set(reflect.ValueOf(docField.Value))
//	}
//
//	return nil
//}

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
