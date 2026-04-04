package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"db/internal/documenterrors"
	"db/internal/services"
	"db/models"
)

type Handler struct {
	services *services.Services
}

func NewHandler(services *services.Services) *Handler {
	return &Handler{services: services}
}

func (h *Handler) health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) putDocument(w http.ResponseWriter, r *http.Request) {
	var dto PutDocumentDto
	if !h.handleDecode(&dto, w, r) {
		return
	}
	doc := models.Document{
		Id:   dto.Document.Id,
		Name: dto.Document.Name,
		Age:  dto.Document.Age,
	}

	id, isCreated, err := h.services.PutDocument(r.Context(), dto.Collection, doc)
	if err != nil {
		h.handleResponseErrors(w, err)
		return
	}
	status := http.StatusOK

	if isCreated {
		status = http.StatusCreated
	}

	res := PutDocumentResponse{
		Document: models.Document{Id: id, Name: doc.Name, Age: doc.Age},
	}

	w.WriteHeader(status)
	if err = json.NewEncoder(w).Encode(res); err != nil {
		slog.Error("Error encoding response:", "error", err)
	}

}

func (h *Handler) getDocument(w http.ResponseWriter, r *http.Request) {
	var dto GetDocumentDto
	if !h.handleDecode(&dto, w, r) {
		return
	}

	doc, err := h.services.GetDocument(r.Context(), dto.Collection, dto.Id)

	if err != nil {
		h.handleResponseErrors(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(GetDocument200Response{
		Document: doc,
	}); err != nil {
		slog.Error("Error encoding response:", "error", err)
	}
}

func (h *Handler) listDocuments(w http.ResponseWriter, r *http.Request) {
	var dto ListDocumentsDto
	if !h.handleDecode(&dto, w, r) {
		return
	}
	docs, err := h.services.ListDocuments(r.Context(), dto.Collection)
	if err != nil {
		h.handleResponseErrors(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(ListDocuments200Response{Documents: docs}); err != nil {
		slog.Error("Error encoding response:", "error", err)
	}
}

func (h *Handler) deleteDocument(w http.ResponseWriter, r *http.Request) {
	var dto DeleteDocumentDto
	if !h.handleDecode(&dto, w, r) {
		return
	}
	if err := h.services.DeleteDocument(r.Context(), dto.Collection, dto.Id); err != nil {
		h.handleResponseErrors(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(DeleteDocument200Response{Id: dto.Id}); err != nil {
		slog.Error("Error encoding response:", "error", err)
	}
}

func (h *Handler) createCollection(w http.ResponseWriter, r *http.Request) {
	var dto CreateCollectionDto
	if !h.handleDecode(&dto, w, r) {
		return
	}

	if err := h.services.CreateCollection(r.Context(), dto.Name); err != nil {
		h.handleResponseErrors(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(CreateCollection201Response{Name: dto.Name}); err != nil {
		slog.Error("Error encoding response:", "error", err)
	}
}

func (h *Handler) listCollections(w http.ResponseWriter, r *http.Request) {
	list, err := h.services.ListCollections(r.Context())
	if err != nil {
		h.handleResponseErrors(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(ListCollections200Response{Collections: list}); err != nil {
		slog.Error("Error encoding response:", "error", err)
	}
}

func (h *Handler) deleteCollection(w http.ResponseWriter, r *http.Request) {
	var dto DeleteCollection200Response
	if !h.handleDecode(&dto, w, r) {
		return
	}
	if err := h.services.DeleteCollection(r.Context(), dto.Name); err != nil {
		h.handleResponseErrors(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(DeleteCollection200Response{Name: dto.Name}); err != nil {
		slog.Error("Error encoding response:", "error", err)
	}

}

func (h *Handler) createIndex(w http.ResponseWriter, r *http.Request) {
	var dto CreateIndexDto
	if !h.handleDecode(&dto, w, r) {
		return
	}

	if err := h.services.CreateIndex(r.Context(), dto.Collection, dto.Name, dto.Fields); err != nil {
		h.handleResponseErrors(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]string{
		"message": "Index created successfully",
	}); err != nil {
		slog.Error("Error encoding response:", "error", err)
	}
}

func (h *Handler) deleteIndex(w http.ResponseWriter, r *http.Request) {
	var dto DeleteIndexDto
	if !h.handleDecode(&dto, w, r) {
		return
	}
	err := h.services.DeleteIndex(r.Context(), dto.Collection, dto.Name)
	if err != nil {
		h.handleResponseErrors(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(map[string]string{
		"message": "Index deleted successfully",
	}); err != nil {
		slog.Error("Error encoding response:", "error", err)
	}
}

func (h *Handler) handleResponseErrors(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, documenterrors.ErrCollectionAlreadyExists):
		http.Error(w, err.Error(), http.StatusConflict)
	case errors.Is(err, documenterrors.ErrCollectionNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, documenterrors.ErrDocumentNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, documenterrors.ErrEmptyCollectionName):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, documenterrors.ErrEmptyPrimaryKey):
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) handleDecode(dto any, w http.ResponseWriter, r *http.Request) bool {
	if err := json.NewDecoder(r.Body).Decode(dto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return false
	}
	return true
}
