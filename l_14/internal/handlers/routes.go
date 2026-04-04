package handlers

import "net/http"

func (h *Handler) InitRoutes() *http.ServeMux {
	m := http.NewServeMux()
	m.HandleFunc("/health", h.health)
	m.HandleFunc("POST /put_document", h.putDocument)
	m.HandleFunc("POST /get_document", h.getDocument)
	m.HandleFunc("POST /list_documents", h.listDocuments)
	m.HandleFunc("POST /delete_document", h.deleteDocument)

	m.HandleFunc("POST /create_collection", h.createCollection)
	m.HandleFunc("POST /list_collections", h.listCollections)
	m.HandleFunc("POST /delete_collection", h.deleteCollection)

	m.HandleFunc("POST /create_index", h.createIndex)
	m.HandleFunc("POST /delete_index", h.deleteIndex)
	return m
}
