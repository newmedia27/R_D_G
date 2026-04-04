package chat

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"project/internal/logger"
	"project/internal/models"
)

type Repository struct {
	collection *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	collection := db.Collection("chats")
	return &Repository{
		collection: collection,
	}
}

func (r *Repository) Create(ctx context.Context, chat *models.Chat) (*models.Chat, error) {
	doc, err := toDocument(chat)
	if err != nil {
		return nil, fmt.Errorf("can't convert to document: %w", err)
	}
	_, err = r.collection.InsertOne(ctx, doc)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, models.ErrChatDuplicate
		}
		return nil, fmt.Errorf("can't insert chat: %w", err)
	}
	return toModel(doc), nil
}

func (r *Repository) FindByID(ctx context.Context, id string) (*models.Chat, error) {
	var doc document
	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	filter := bson.M{"_id": ID}
	err = r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, models.ErrChatNotFound
		}
		return nil, fmt.Errorf("can't find chat: %w", err)
	}
	chat := toModel(&doc)
	return chat, nil
}

func (r *Repository) FindAllByUserID(ctx context.Context, userID string) ([]*models.Chat, error) {
	var docs []document
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, models.ErrInvalidUserIds
	}
	filter := bson.M{"members": id}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list chats: %w", err)
	}
	defer func(c context.Context) {
		if err = cursor.Close(c); err != nil {
			logger.LogErrorContext(ctx, err)
		}
	}(ctx)
	if err = cursor.All(ctx, &docs); err != nil {
		return nil, fmt.Errorf("failed to list chats: %w", err)
	}
	res := make([]*models.Chat, 0, len(docs))
	for _, doc := range docs {
		res = append(res, toModel(&doc))
	}
	return res, nil
}

func (r *Repository) FindPrivateChat(ctx context.Context, id1, id2 string) (*models.Chat, error) {
	usr1, err := primitive.ObjectIDFromHex(id1)
	if err != nil {
		return nil, models.ErrInvalidUserIds
	}
	usr2, err := primitive.ObjectIDFromHex(id2)
	if err != nil {
		return nil, models.ErrInvalidUserIds
	}
	filter := bson.M{
		"type": models.ChatTypePrivate,
		"members": bson.M{
			"$all":  bson.A{usr1, usr2},
			"$size": 2,
		},
	}
	var doc document

	err = r.collection.FindOne(ctx, filter).Decode(&doc)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, models.ErrChatNotFound
		}
		return nil, fmt.Errorf("can't find chat: %w", err)
	}
	return toModel(&doc), nil
}

func (r *Repository) AddMember(ctx context.Context, chatID, userID string) error {
	chatId, err := primitive.ObjectIDFromHex(chatID)
	if err != nil {
		return fmt.Errorf("invalid chat id: %w", err)
	}
	usrId, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("invalid user id: %w", err)
	}
	filter := bson.M{"_id": chatId}
	update := bson.M{"$addToSet": bson.M{"members": usrId}}
	_, err = r.collection.UpdateOne(ctx, filter, update)

	if err != nil {
		return fmt.Errorf("can't add member: %w", err)
	}
	return nil
}

func (r *Repository) RemoveMember(ctx context.Context, chatID, userID string) error {
	chatId, err := primitive.ObjectIDFromHex(chatID)
	if err != nil {
		return fmt.Errorf("invalid chat id: %w", err)
	}
	usrId, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("invalid user id: %w", err)
	}
	filter := bson.M{"_id": chatId}
	update := bson.M{"$pull": bson.M{"members": usrId}}
	_, err = r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("can't remove member: %w", err)
	}
	return nil
}

func (r *Repository) UpdateLastMessage(ctx context.Context, chatID string, msg *models.Message) error {
	chatId, err := primitive.ObjectIDFromHex(chatID)
	if err != nil {
		return fmt.Errorf("invalid chat id: %w", err)
	}
	usrId, err := primitive.ObjectIDFromHex(msg.UserID)
	if err != nil {
		return fmt.Errorf("invalid user id: %w", err)
	}
	filter := bson.M{"_id": chatId}
	update := bson.M{"$set": bson.M{
		"last_message": lastMessageDocument{
			Text:      msg.Text,
			UserID:    usrId,
			CreatedAt: msg.CreatedAt,
		},
		"updated_at": msg.UpdatedAt,
	}}
	_, err = r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("can't update last message: %w", err)
	}
	return nil
}

func (r *Repository) FindByOwnerAndName(ctx context.Context, ownerID, name string) (*models.Chat, error) {
	oID, err := primitive.ObjectIDFromHex(ownerID)
	if err != nil {
		return nil, fmt.Errorf("invalid owner id: %w", err)
	}
	var doc document
	filter := bson.M{
		"type":     models.ChatTypeGroup,
		"owner_id": oID,
		"name":     name,
	}
	err = r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, models.ErrChatNotFound
		}
		return nil, fmt.Errorf("find group chat: %w", err)
	}
	return toModel(&doc), nil
}
