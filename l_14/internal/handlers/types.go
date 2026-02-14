package handlers

import "db/models"

type PutDocumentDto struct {
	Collection string          `json:"collection_name"`
	Document   models.Document `json:"document"`
}
type PutDocumentResponse struct {
	Document models.Document `json:"document"`
}
type GetDocumentDto struct {
	Collection string `json:"collection_name"`
	Id         string `json:"id"`
}
type GetDocument200Response struct {
	Document models.Document `json:"document"`
}
type ListDocumentsDto struct {
	Collection string `json:"collection_name"`
}
type ListDocuments200Response struct {
	Documents []models.Document `json:"documents"`
}
type DeleteDocumentDto struct {
	Collection string `json:"collection_name"`
	Id         string `json:"id"`
}
type DeleteDocument200Response struct {
	Id string `json:"id"`
}
type CreateCollectionDto struct {
	Name string `json:"name"`
}
type CreateCollection201Response struct {
	Name string `json:"name"`
}
type ListCollections200Response struct {
	Collections []string `json:"collections"`
}
type DeleteCollectionDto struct {
	Name string `json:"name"`
}
type DeleteCollection200Response struct {
	Name string `json:"name"`
}

// Value asc/desc
type CreateIndexDto struct {
	Collection string              `json:"collection_name"`
	Name       string              `json:"index_name"`
	Fields     []models.IndexField `json:"fields"`
}
type DeleteIndexDto struct {
	Collection string `json:"collection_name"`
	Name       string `json:"index_name"`
}
