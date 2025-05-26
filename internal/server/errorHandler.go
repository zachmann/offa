package server

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

func handleError(ctx *fiber.Ctx, err error) error {
	// Status code defaults to 500
	code := fiber.StatusInternalServerError
	msg := err.Error()

	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
		msg = e.Error()
	}
	return ctx.Status(code).JSON(fiber.Map{"error": msg})
}
