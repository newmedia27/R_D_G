package main

import (
	"fmt"
	"log/slog"

	"github.com/newmedia27/R_D_G/l_3/pkg/documentstore"
	"github.com/newmedia27/R_D_G/l_3/pkg/printobject"
)

var docValid = &documentstore.Document{
	Fields: map[string]documentstore.DocumentField{
		"key": {
			Type:  documentstore.DocumentFieldTypeString,
			Value: "key_1",
		},
		"name": {
			Type:  documentstore.DocumentFieldTypeString,
			Value: "Sviatoslav",
		},
	},
}
var docValid2 = &documentstore.Document{
	Fields: map[string]documentstore.DocumentField{
		"key": {
			Type:  documentstore.DocumentFieldTypeString,
			Value: "key_2",
		},
		"name": {
			Type:  documentstore.DocumentFieldTypeString,
			Value: "Some name",
		},
	},
}
var docValid3 = &documentstore.Document{
	Fields: map[string]documentstore.DocumentField{
		"key": {
			Type:  documentstore.DocumentFieldTypeString,
			Value: "key_3",
		},
		"name": {
			Type:  documentstore.DocumentFieldTypeString,
			Value: "Some name",
		},
	},
}

var docInvalid = &documentstore.Document{
	Fields: map[string]documentstore.DocumentField{

		"name": {
			Type:  documentstore.DocumentFieldTypeString,
			Value: "invalid document",
		},
	},
}

func main() {
	var err error

	err = documentstore.Put(docValid)
	if err != nil {
		slog.Error("Failed to put document",
			slog.Any("doc", docValid),
			slog.Any("error", err),
		)
	}
	err = documentstore.Put(docValid2)
	if err != nil {
		slog.Error("Failed to put document",
			slog.Any("doc", docValid2),
			slog.Any("error", err),
		)
	}
	err = documentstore.Put(docInvalid)
	if err != nil {
		slog.Error("Failed to put document",
			slog.Any("doc", docInvalid),
			slog.Any("error", err),
		)
	}
	err = documentstore.Put(docValid3)
	if err != nil {
		slog.Error("Failed to put document",
			slog.Any("doc", docInvalid),
			slog.Any("error", err),
		)
	}

	doc, ok := documentstore.Get("key_1")
	if !ok {
		fmt.Println("Document not found")
	}

	printobject.PrintObject("Document: ", doc)
	isDeleted := documentstore.Delete("key_1")
	fmt.Println("Document deleted:", isDeleted)

	docs := documentstore.List()
	printobject.PrintObject("List: ", docs)
}
