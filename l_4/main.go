package main

import (
	"fmt"

	"github.com/newmedia27/R_D_G/l_4/pkg/documentstore"
	"github.com/newmedia27/R_D_G/l_4/pkg/printobject"
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
	store := documentstore.NewStore()
	keyCfg := &documentstore.CollectionConfig{PrimaryKey: "key"}
	var ok bool
	var keyCollection *documentstore.Collection
	var doc *documentstore.Document

	ok, keyCollection = store.CreateCollection("key", keyCfg)
	if !ok {
		fmt.Println("Collection already exists")
	}

	keyCollection.Put(*docValid)
	keyCollection.Put(*docValid2)
	keyCollection.Put(*docValid3)

	keyCollection, ok = store.GetCollection("key")
	if ok {
		printobject.PrintObject("getKeyCol", keyCollection)
		printobject.PrintObject("Store", store)
	}

	if doc, ok = keyCollection.Get("key_1"); ok {
		printobject.PrintObject("GetDoc", doc)
	}
	if deletedDoc := keyCollection.Delete("key_1"); deletedDoc {
		fmt.Println("Deleted doc: key_1")
		printobject.PrintObject("Store", store)
	}

	docList := keyCollection.List()
	printobject.PrintObject("List", docList)

	deleted := store.DeleteCollection("key")
	if deleted {
		printobject.PrintObject("DeleteKeyCollection", store)
	}

}
