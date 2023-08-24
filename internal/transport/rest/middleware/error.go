package middleware

import (
	"errors"
	"github.com/ensiouel/apperror"
	"github.com/gofiber/fiber/v2"
)

func Error() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		var e *fiber.Error
		if errors.As(err, &e) {
			return c.Status(e.Code).JSON(fiber.Map{"error": e.Error()})
		}

		var apperr *apperror.Error
		if errors.As(err, &apperr) {
			switch apperr.Code {
			case apperror.Internal.Code:
				c.Status(fiber.StatusInternalServerError)
			case apperror.NotFound.Code:
				c.Status(fiber.StatusNotFound)
			case apperror.AlreadyExists.Code, apperror.BadRequest.Code:
				c.Status(fiber.StatusBadRequest)
			case apperror.Unauthorized.Code:
				c.Status(fiber.StatusUnauthorized)
			}

			return c.JSON(fiber.Map{"error": apperr})
		}

		return c.Status(fiber.StatusTeapot).JSON(fiber.Map{"error": apperror.Unknown.WithError(err)})
	}
}
