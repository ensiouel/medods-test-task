package dto

import (
	"github.com/ensiouel/apperror"
	"github.com/google/uuid"
)

type AuthRequest struct {
	GUID uuid.UUID `form:"guid"`
}

func (request AuthRequest) Validate() error {
	if request.GUID == uuid.Nil {
		return apperror.BadRequest.WithMessage("invalid guid")
	}

	return nil
}

type RefreshRequest struct {
	SessionID    uuid.UUID `params:"session_id"`
	RefreshToken []byte    `json:"refresh_token"`
}

func (request RefreshRequest) Validate() error {
	if request.SessionID == uuid.Nil {
		return apperror.BadRequest.WithMessage("invalid session id")
	}

	if len(request.RefreshToken) == 0 {
		return apperror.BadRequest.WithMessage("invalid refresh token")
	}

	return nil
}
