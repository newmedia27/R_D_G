package message

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
)

type Repository struct {
	collection *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	collection := db.Collection("messages")
	return &Repository{
		collection: collection,
	}
}

func (r *Repository) Create(ctx context.Context, msg *models.Message) (*models.Message, error) {
	doc, err := toDocument(msg)
	if err != nil {
		return nil, fmt.Errorf("can't convert to document: %w", err)
	}
	_, err = r.collection.InsertOne(ctx, doc)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, models.ErrMessageDuplicate
		}
		return nil, fmt.Errorf("can't insert message: %w", err)
	}
	return toModel(doc), nil
}

func (r *Repository) FindByID(ctx context.Context, id string) (*models.Message, error) {
	var doc document
	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	filter := bson.M{"_id": ID}
	err = r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, models.ErrMessageNotFound
		}
		return nil, fmt.Errorf("can't find message: %w", err)
	}
	return toModel(&doc), nil
}

func (r *Repository) FindByChatID(ctx context.Context, chatID, lastSeeMsgId string, limit int) ([]*models.Message, error) {
	var docs []document
	chatId, err := primitive.ObjectIDFromHex(chatID)
	if err != nil {
		return nil, fmt.Errorf("invalid chat id: %w", err)
	}
	filter := bson.M{"chat_id": chatId}
	if lastSeeMsgId != "" {
		before, err := primitive.ObjectIDFromHex(lastSeeMsgId)
		if err != nil {
			return nil, fmt.Errorf("invalid last see message id: %w", err)
		}
		filter["_id"] = bson.M{"$lt": before}
	}
	opts := options.Find().SetSort(bson.D{{"created_at", -1}}).SetLimit(int64(limit))
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list messages: %w", err)
	}
	defer func(c context.Context) {
		if err = cursor.Close(c); err != nil {
			logger.LogErrorContext(c, err)
		}
	}(ctx)
	if err = cursor.All(ctx, &docs); err != nil {
		return nil, fmt.Errorf("failed to list messages: %w", err)
	}
	res := make([]*models.Message, 0, len(docs))
	for _, doc := range docs {
		res = append(res, toModel(&doc))
	}
	for i, j := 0, len(docs)-1; i < j; i, j = i+1, j-1 {
		res[i], res[j] = res[j], res[i]
	}
	return res, nil
}

func (r *Repository) UpdateMessage(ctx context.Context, msg *models.Message) error {
	ID, err := primitive.ObjectIDFromHex(msg.ID)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	filter := bson.M{"_id": ID}
	update := bson.M{"$set": bson.M{
		"text":       msg.Text,
		"updated_at": msg.UpdatedAt,
		"is_edited":  true,
	}}
	_, err = r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("can't update message: %w", err)
	}
	return nil
}

func (r *Repository) SoftDelete(ctx context.Context, msgID string) error {
	ID, err := primitive.ObjectIDFromHex(msgID)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	filter := bson.M{"_id": ID}
	update := bson.M{"$set": bson.M{
		"is_deleted": true,
	}}
	_, err = r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("can't soft delete message: %w", err)
	}
	return nil
}

func (r *Repository) GetMessage(ctx context.Context, id string) (*models.Message, error) {
	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	var doc document
	err = r.collection.FindOne(ctx, bson.M{"_id": ID}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, models.ErrMessageNotFound
		}
		return nil, fmt.Errorf("find message: %w", err)
	}
	return toModel(&doc), nil
}
