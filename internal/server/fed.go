package server

import (
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"github.com/zachmann/go-oidfed/pkg/constants"
)

func addFederationEndpoints(s fiber.Router) {
	s.Get("/.well-known/openid-federation", handleEntityConfiguration)
}

func handleEntityConfiguration(ctx *fiber.Ctx) error {
	jwt, err := federationLeafEntity.EntityConfigurationJWT()
	if err != nil {
		log.WithError(err).Error("Failed to get entity configuration")
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get entity configuration")
	}
	ctx.Set(fiber.HeaderContentType, constants.ContentTypeEntityStatement)
	return ctx.Send(jwt)
}
