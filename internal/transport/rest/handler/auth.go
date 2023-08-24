package handler

import (
	"github.com/ensiouel/apperror"
	"github.com/gofiber/fiber/v2"
	"medods-test-task/internal/dto"
	"medods-test-task/internal/model"
	"medods-test-task/internal/service"
)

type SessionHandler struct {
	sessionService service.SessionService
}

func NewSessionHandler(sessionService service.SessionService) *SessionHandler {
	return &SessionHandler{
		sessionService: sessionService,
	}
}

func (handler *SessionHandler) Register(router fiber.Router) {
	router.Get("/auth", handler.auth)
	router.Post("/:session_id/refresh", handler.refresh)
}

func (handler *SessionHandler) auth(c *fiber.Ctx) error {
	var authRequest dto.AuthRequest
	err := c.QueryParser(&authRequest)
	if err != nil {
		return apperror.BadRequest.WithError(err)
	}

	err = authRequest.Validate()
	if err != nil {
		return err
	}

	var tokenPair model.TokenPair
	tokenPair, err = handler.sessionService.Create(c.Context(), authRequest.GUID)
	if err != nil {
		return err
	}

	return c.JSON(tokenPair)
}

func (handler *SessionHandler) refresh(c *fiber.Ctx) error {
	var refreshRequest dto.RefreshRequest

	err := c.BodyParser(&refreshRequest)
	if err != nil {
		return apperror.BadRequest.WithError(err)
	}

	err = c.ParamsParser(&refreshRequest)
	if err != nil {
		return apperror.BadRequest.WithError(err)
	}

	err = refreshRequest.Validate()
	if err != nil {
		return err
	}

	var tokenPair model.TokenPair
	tokenPair, err = handler.sessionService.Update(c.Context(), refreshRequest.SessionID, refreshRequest.RefreshToken)
	if err != nil {
		return err
	}

	return c.JSON(tokenPair)
}
