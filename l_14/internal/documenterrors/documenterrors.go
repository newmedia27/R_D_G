package documenterrors

import "errors"

var (
	ErrCollectionAlreadyExists = errors.New("collection already exists")
	ErrCollectionNotFound      = errors.New("collection not found")
	ErrDocumentNotFound        = errors.New("document not found")
	ErrEmptyCollectionName     = errors.New("empty collection name")
	ErrEmptyPrimaryKey         = errors.New("empty primary key")
)
