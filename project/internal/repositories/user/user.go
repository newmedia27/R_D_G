package user

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"project/internal/logger"
	"project/internal/models"
	"project/internal/repositories/common"
)

type Repository struct {
	collection *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	collection := db.Collection("users")
	return &Repository{
		collection: collection,
	}
}

func (r *Repository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var u document
	filter := bson.M{"email": email}
	err := r.collection.FindOne(ctx, filter).Decode(&u)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, models.ErrUserNotFound
		}
		return nil, fmt.Errorf("can't find user: %w", err)
	}
	user := toModel(&u)
	return user, nil
}

func (r *Repository) FindByID(ctx context.Context, id string) (*models.User, error) {
	var u document
	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	filter := bson.M{"_id": ID}
	err = r.collection.FindOne(ctx, filter).Decode(&u)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, models.ErrUserNotFound
		}
		return nil, fmt.Errorf("can't find user: %w", err)
	}
	user := toModel(&u)
	return user, nil
}

func (r *Repository) Create(ctx context.Context, u *models.User) (*models.User, error) {
	doc, err := toDocument(u)
	if err != nil {
		return nil, fmt.Errorf("can't convert to document: %w", err)
	}
	_, err = r.collection.InsertOne(ctx, doc)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, models.ErrUserDuplicate
		}
		return nil, fmt.Errorf("can't insert user: %w", err)
	}
	return toModel(doc), nil
}

func (r *Repository) FindByIds(ctx context.Context, ids []string) ([]*models.User, error) {
	userIds := make([]primitive.ObjectID, len(ids))

	for i, id := range ids {
		userId, err := common.ParseOrGeneratePrimitiveId(id)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return nil, fmt.Errorf("invalid id: %w", err)
			}
			return nil, fmt.Errorf("invalid id: %w", err)
		}
		userIds[i] = userId
	}

	var docs []*document
	opts := options.Find().SetSort(bson.D{{Key: "_id", Value: 1}})
	filter := bson.M{"_id": bson.M{"$in": userIds}}
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("can't find users: %w", err)
	}
	defer func(c context.Context) {
		if err = cursor.Close(c); err != nil {
			logger.LogErrorContext(c, err)
		}
	}(ctx)
	if err = cursor.All(ctx, &docs); err != nil {
		return nil, fmt.Errorf("can't find users: %w", err)
	}
	res := make([]*models.User, 0, len(docs))
	for _, doc := range docs {
		res = append(res, toModel(doc))
	}
	return res, nil
}

func (r *Repository) Search(ctx context.Context, userId, query string) ([]*models.User, error) {
	id, err := common.ParseOrGeneratePrimitiveId(userId)
	if err != nil {
		return nil, err
	}
	filter := bson.M{
		"_id": bson.M{"$ne": id},
		"$or": bson.A{
			bson.M{"username": bson.M{"$regex": query, "$options": "i"}},
			bson.M{"email": bson.M{"$regex": query, "$options": "i"}},
		},
	}
	opts := options.Find().SetLimit(50).SetSort(bson.D{{Key: "username", Value: 1}})

	var docs []*document
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("can't find users: %w", err)
	}
	defer func(c context.Context) {
		if err = cursor.Close(c); err != nil {
			logger.LogErrorContext(c, err)
		}
	}(ctx)
	if err = cursor.All(ctx, &docs); err != nil {
		return nil, fmt.Errorf("can't find users: %w", err)
	}
	res := make([]*models.User, 0, len(docs))
	for _, doc := range docs {
		res = append(res, toModel(doc))
	}
	return res, nil
}
