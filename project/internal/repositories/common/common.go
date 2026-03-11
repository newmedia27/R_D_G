package common

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"project/internal/models"
)

func ParseOrGeneratePrimitiveId(id string) (primitive.ObjectID, error) {
	if id == "" {
		return primitive.NewObjectID(), nil
	}
	Id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.ObjectID{}, models.ErrIncorrectId
	}
	return Id, nil
}
