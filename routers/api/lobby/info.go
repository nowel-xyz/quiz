package lobby

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nowel-xyz/quiz/database/models"

	"github.com/nowel-xyz/quiz/service/lobby"
	
)

func SetupLobbyInfoRoutes(router fiber.Router) {
	router.Get("/:id", func(c *fiber.Ctx) error {
		user := c.Locals("user").(models.User)
		lobbyID := c.Params("id")

		lobby, err := service_lobby.GetLobbyData(c, lobbyID, user)
		if err != nil {
			return err
		}

		return c.JSON(lobby)
	})
}
