package session

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
	collection := db.Collection("sessions")
	return &Repository{
		collection: collection,
	}
}

func (r *Repository) Create(ctx context.Context, s *models.Session) error {
	doc, err := toDocument(s)
	if err != nil {
		return fmt.Errorf("cant convert to document: %w", err)
	}
	_, err = r.collection.InsertOne(ctx, doc)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return models.ErrSessionDuplicate
		}
		return fmt.Errorf("can't create session: %w", err) //TODO:
	}
	return nil
}

func (r *Repository) FindByToken(ctx context.Context, token string) (*models.Session, error) {
	var doc document
	filter := bson.M{"refresh_token_hash": token}
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, models.ErrSessionNotFound
		}
		return nil, fmt.Errorf("can't find session: %w", err)
	}
	return toModel(&doc), nil
}
func (r *Repository) FindByUserID(ctx context.Context, userID string) (*models.Session, error) {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}
	var doc document
	filter := bson.M{"user_id": id}
	err = r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, models.ErrSessionNotFound
		}
		return nil, fmt.Errorf("can't find session: %w", err)
	}
	return toModel(&doc), nil
}

func (r *Repository) FindByDeviceId(ctx context.Context, deviceID string) (*models.Session, error) {
	var doc document
	filter := bson.M{"device_session_id": deviceID}
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, models.ErrSessionNotFound
		}
		return nil, fmt.Errorf("can't find session: %w", err)
	}
	return toModel(&doc), nil
}

func (r *Repository) DeleteByDeviceId(ctx context.Context, deviceID string) error {
	filter := bson.M{"device_session_id": deviceID}
	_, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("can't delete session: %w", err)
	}
	return nil
}

func (r *Repository) DeleteAllByUserID(ctx context.Context, userID string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("invalid user id: %w", err)
	}
	filter := bson.M{"user_id": id}
	_, err = r.collection.DeleteMany(ctx, filter)
	if err != nil {
		return fmt.Errorf("can't delete session: %w", err)
	}
	return nil
}

func (r *Repository) FindAllByUserID(ctx context.Context, userID string) ([]*models.Session, error) {
	var docs []document
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}
	filter := bson.M{"user_id": id}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("can't find sessions: %w", err)
	}
	defer func(c context.Context) {
		if err = cursor.Close(c); err != nil {
			logger.LogErrorContext(c, err)
		}
	}(ctx)

	err = cursor.All(ctx, &docs)
	if err != nil {
		return nil, fmt.Errorf("can't find sessions: %w", err)
	}
	sessions := make([]*models.Session, 0, len(docs))
	for _, doc := range docs {
		sessions = append(sessions, toModel(&doc))
	}
	return sessions, nil
}

func (r *Repository) Update(ctx context.Context, s *models.Session) error {
	id, err := primitive.ObjectIDFromHex(s.ID)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{
		"refresh_token_hash": s.RefreshTokenHash,
		"last_used_at":       s.LastUsedAt,
		"user_agent":         s.UserAgent,
		"ip":                 s.IP,
		"fingerprint":        s.FingerPrint,
	}}
	_, err = r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("can't update session: %w", err)
	}
	return nil
}

func (r *Repository) FindByUserIDAndFingerprint(ctx context.Context, userID, fingerprint string) (*models.Session, error) {
	var doc document
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}
	filter := bson.M{"user_id": id, "fingerprint": fingerprint}
	err = r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, models.ErrSessionNotFound
		}
	}
	return toModel(&doc), nil
}
