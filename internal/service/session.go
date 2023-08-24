package service

import (
	"context"
	"github.com/ensiouel/apperror"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"medods-test-task/internal/config"
	"medods-test-task/internal/model"
	"medods-test-task/internal/storage"
	"time"
)

type SessionService interface {
	Create(ctx context.Context, guid uuid.UUID) (model.TokenPair, error)
	Update(ctx context.Context, sessionID uuid.UUID, refreshToken []byte) (model.TokenPair, error)
}

type SessionServiceImpl struct {
	storage storage.SessionStorage
	conf    config.Session
}

func NewSessionService(storage storage.SessionStorage, conf config.Session) *SessionServiceImpl {
	return &SessionServiceImpl{
		storage: storage,
		conf:    conf,
	}
}

func (service *SessionServiceImpl) Create(ctx context.Context, guid uuid.UUID) (model.TokenPair, error) {
	now := time.Now()

	session := model.Session{
		ID:        uuid.New(),
		GUID:      guid,
		CreatedAt: now,
		UpdatedAt: now,
	}

	tokenPair, err := service.signTokenPair(session.ID, guid)
	if err != nil {
		return model.TokenPair{}, err
	}

	session.HashedRefreshToken, err = service.hashFromToken(tokenPair.RefreshToken)
	if err != nil {
		return model.TokenPair{}, err
	}

	err = service.storage.Create(ctx, session)
	if err != nil {
		return model.TokenPair{}, err
	}

	return tokenPair, nil
}

func (service *SessionServiceImpl) Update(ctx context.Context, sessionID uuid.UUID, refreshToken []byte) (model.TokenPair, error) {
	session, err := service.storage.GetByID(ctx, sessionID)
	if err != nil {
		return model.TokenPair{}, err
	}

	err = bcrypt.CompareHashAndPassword(session.HashedRefreshToken, refreshToken)
	if err != nil {
		return model.TokenPair{}, apperror.BadRequest.WithMessage("invalid refresh token")
	}

	now := time.Now()

	if session.UpdatedAt.Sub(now) >= service.conf.RefreshTokenExpiration {
		err = service.storage.DeleteByID(ctx, session.ID)
		if err != nil {
			return model.TokenPair{}, err
		}

		return model.TokenPair{}, apperror.BadRequest.WithMessage("refresh token expired")
	}

	var tokenPair model.TokenPair
	tokenPair, err = service.signTokenPair(session.ID, session.GUID)
	if err != nil {
		return model.TokenPair{}, err
	}

	session.UpdatedAt = now

	session.HashedRefreshToken, err = service.hashFromToken(tokenPair.RefreshToken)
	if err != nil {
		return model.TokenPair{}, err
	}

	err = service.storage.Update(ctx, session)
	if err != nil {
		return model.TokenPair{}, err
	}

	return tokenPair, nil
}

func (service *SessionServiceImpl) signTokenPair(sessionID uuid.UUID, guid uuid.UUID) (model.TokenPair, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, model.AccessTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(service.conf.AccessTokenExpiration)),
		},
		SessionID: sessionID,
		GUID:      guid,
	})

	signedAccessToken, err := accessToken.SignedString(service.conf.Secret)
	if err != nil {
		return model.TokenPair{}, apperror.Internal.WithError(err)
	}

	refreshToken := uuid.New()

	tokenPair := model.TokenPair{
		AccessToken:  signedAccessToken,
		RefreshToken: refreshToken[:],
	}

	return tokenPair, nil
}

func (service *SessionServiceImpl) hashFromToken(token []byte) ([]byte, error) {
	bytes, err := bcrypt.GenerateFromPassword(token, bcrypt.DefaultCost)
	if err != nil {
		return nil, apperror.Internal.WithError(err)
	}

	return bytes, nil
}
