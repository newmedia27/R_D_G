package services

import (
	"context"
	"fmt"

	"db/internal/repositories"
	"db/models"
	"go.mongodb.org/mongo-driver/bson"
)

type Services struct {
	repo *repositories.Repositories
}

func NewServices(r *repositories.Repositories) *Services {
	return &Services{
		repo: r,
	}
}

func (s *Services) PutDocument(ctx context.Context, collectionName string, document models.Document) (string, bool, error) {
	return s.repo.PutDocument(ctx, collectionName, document)
}

func (s *Services) GetDocument(ctx context.Context, collectionName string, id string) (models.Document, error) {
	return s.repo.GetDocument(ctx, collectionName, id)
}

func (s *Services) ListDocuments(ctx context.Context, collectionName string) ([]models.Document, error) {
	return s.repo.ListDocuments(ctx, collectionName)
}

func (s *Services) DeleteDocument(ctx context.Context, collectionName string, id string) error {
	return s.repo.DeleteDocument(ctx, collectionName, id)
}

func (s *Services) CreateCollection(ctx context.Context, name string) error {
	return s.repo.CreateCollection(ctx, name)
}

func (s *Services) ListCollections(ctx context.Context) ([]string, error) {
	return s.repo.ListCollections(ctx)
}

func (s *Services) DeleteCollection(ctx context.Context, name string) error {
	return s.repo.DeleteCollection(ctx, name)
}

func (s *Services) CreateIndex(ctx context.Context, collectionName, indexName string, fields []models.IndexField) error {

	b := bson.D{}

	for _, field := range fields {
		switch field.Value {
		case "asc":
			b = append(b, bson.E{Key: field.Name, Value: 1})
		case "desc":
			b = append(b, bson.E{Key: field.Name, Value: -1})
		default:
			return fmt.Errorf("unsupported index field: %s", field.Value)

		}
	}

	return s.repo.CreateIndex(ctx, collectionName, indexName, b)
}

func (s *Services) DeleteIndex(ctx context.Context, collectionName, indexName string) error {
	return s.repo.DeleteIndex(ctx, collectionName, indexName)
}
