package repositories

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"db/internal/documenterrors"
	"db/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repositories struct {
	db *mongo.Database
}

func NewRepositories(db *mongo.Database) *Repositories {
	return &Repositories{
		db,
	}
}

func (r *Repositories) PutDocument(ctx context.Context, collectionName string, document models.Document) (string, bool, error) {

	if document.Id == "" {
		document.Id = primitive.NewObjectID().Hex()
	}

	Id, err := primitive.ObjectIDFromHex(document.Id)
	if err != nil {
		return "", false, fmt.Errorf("invalid id: %w", err)
	}

	filter := bson.M{
		"_id": Id,
	}
	update := bson.M{
		"$set": document,
	}
	opt := options.Update().SetUpsert(true)

	result, err := r.db.Collection(collectionName).UpdateOne(ctx, filter, update, opt)
	if err != nil {
		return "", false, fmt.Errorf("failed to update document: %w", err)
	}

	if result.UpsertedID != nil {
		return result.UpsertedID.(primitive.ObjectID).Hex(), true, nil
	}
	return document.Id, false, nil
}

func (r *Repositories) GetDocument(ctx context.Context, collectionName string, id string) (models.Document, error) {
	Id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Document{}, fmt.Errorf("invalid id: %w", documenterrors.ErrDocumentNotFound)
	}

	var doc models.Document
	filter := bson.M{"_id": Id}
	err = r.db.Collection(collectionName).FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.Document{}, documenterrors.ErrDocumentNotFound
		}
		return models.Document{}, fmt.Errorf("failed to get document: %w", err)
	}
	return doc, nil
}

func (r *Repositories) ListDocuments(ctx context.Context, collectionName string) ([]models.Document, error) {
	var docs []models.Document
	filter := bson.D{}

	cursor, err := r.db.Collection(collectionName).Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list documents: %w", err)
	}
	defer func(c context.Context) {
		err = cursor.Close(c)
		slog.Error("failed to close cursor", "error", err)
	}(ctx)
	if err = cursor.All(ctx, &docs); err != nil {
		return nil, fmt.Errorf("failed to list documents: %w", err)
	}

	return docs, nil
}

func (r *Repositories) DeleteDocument(ctx context.Context, collectionName string, id string) error {
	Id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", documenterrors.ErrDocumentNotFound)
	}
	filter := bson.M{"_id": Id}
	_, err = r.db.Collection(collectionName).DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete document: %w", documenterrors.ErrDocumentNotFound)
	}
	return nil
}

func (r *Repositories) CreateCollection(ctx context.Context, name string) error {
	filter := bson.M{"name": name}
	collections, err := r.db.ListCollectionNames(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to list collections: %w", err)
	}

	if len(collections) > 0 {
		return fmt.Errorf("CreateCollection: %w ", documenterrors.ErrCollectionAlreadyExists)
	}

	if err = r.db.CreateCollection(ctx, name); err != nil {
		return fmt.Errorf("failed to create collection: %w", err)
	}
	return nil
}

func (r *Repositories) ListCollections(ctx context.Context) ([]string, error) {
	filter := bson.D{}
	list, err := r.db.ListCollectionNames(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list collections: %w", err)
	}
	return list, nil
}

func (r *Repositories) DeleteCollection(ctx context.Context, name string) error {
	filter := bson.M{"name": name}

	collections, err := r.db.ListCollectionNames(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to list collections: %w", err)
	}

	if len(collections) == 0 {
		return fmt.Errorf("DeleteCollection: %w ", documenterrors.ErrCollectionNotFound)
	}

	if err = r.db.Collection(name).Drop(ctx); err != nil {
		return fmt.Errorf("failed to drop collection: %w", err)
	}
	return nil
}

func (r *Repositories) CreateIndex(ctx context.Context, collectionName, indexName string, fields bson.D) error {
	index := mongo.IndexModel{
		Keys:    fields,
		Options: options.Index().SetName(indexName),
	}
	_, err := r.db.Collection(collectionName).Indexes().CreateOne(ctx, index)
	if err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}
	return nil
}

func (r *Repositories) DeleteIndex(ctx context.Context, collectionName, indexName string) error {
	_, err := r.db.Collection(collectionName).Indexes().DropOne(ctx, indexName)
	if err != nil {
		return fmt.Errorf("failed to drop index: %w", err)
	}
	return nil
}
