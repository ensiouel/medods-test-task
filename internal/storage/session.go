package storage

import (
	"context"
	"github.com/ensiouel/apperror"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"medods-test-task/internal/model"
)

type SessionStorage interface {
	Create(ctx context.Context, session model.Session) error
	GetByID(ctx context.Context, sessionID uuid.UUID) (model.Session, error)
	Update(ctx context.Context, session model.Session) error
	DeleteByID(ctx context.Context, sessionID uuid.UUID) error
}

type SessionStorageImpl struct {
	collection *mongo.Collection
}

func NewSessionStorage(client *mongo.Client) *SessionStorageImpl {
	return &SessionStorageImpl{
		collection: client.Database("mongo").Collection("session"),
	}
}

func (storage *SessionStorageImpl) Create(ctx context.Context, session model.Session) error {
	_, err := storage.collection.InsertOne(ctx, session)
	if err != nil {
		return apperror.Internal.WithError(err)
	}

	return nil
}

func (storage *SessionStorageImpl) GetByID(ctx context.Context, sessionID uuid.UUID) (model.Session, error) {
	filter := bson.D{{"id", sessionID}}

	var session model.Session
	err := storage.collection.FindOne(ctx, filter).Decode(&session)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return model.Session{}, apperror.NotFound.WithError(err)
		}

		return model.Session{}, apperror.Internal.WithError(err)
	}

	return session, nil
}

func (storage *SessionStorageImpl) Update(ctx context.Context, session model.Session) error {
	filter := bson.D{{"id", session.ID}}
	update := bson.D{{"$set", bson.D{
		{"hashed_refresh_token", session.HashedRefreshToken},
		{"updated_at", session.UpdatedAt},
	}}}

	_, err := storage.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return apperror.Internal.WithError(err)
	}

	return nil
}

func (storage *SessionStorageImpl) DeleteByID(ctx context.Context, sessionID uuid.UUID) error {
	filter := bson.D{{"id", sessionID}}

	_, err := storage.collection.DeleteOne(ctx, filter)
	if err != nil {
		return apperror.Internal.WithError(err)
	}

	return nil
}
