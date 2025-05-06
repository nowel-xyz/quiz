package lobby

import (
	"github.com/gofiber/fiber/v2"
)

func SetupLobbyRoutes(router fiber.Router) {
	SetupLobbyCreateRoutes(router)
}
