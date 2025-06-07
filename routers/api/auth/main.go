package auth

import (
	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(router fiber.Router) {
	auth := router.Group("/auth")
	SetupLoginRoutes(auth)
	SetupLogoutRoutes(auth)
	SetupRegisterRoutes(auth)
}
