package model

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

type Session struct {
	ID                 uuid.UUID `bson:"id" json:"id"`
	GUID               uuid.UUID `bson:"guid" json:"guid"`
	HashedRefreshToken []byte    `bson:"hashed_refresh_token" json:"hashed_refresh_token"`
	CreatedAt          time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt          time.Time `bson:"updated_at" json:"updated_at"`
}

type AccessTokenClaims struct {
	jwt.RegisteredClaims
	SessionID uuid.UUID `json:"session_id"`
	GUID      uuid.UUID `json:"guid"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken []byte `json:"refresh_token"`
}
