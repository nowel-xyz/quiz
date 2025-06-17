package lobby

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/nowel-xyz/quiz/database/models"
	"github.com/nowel-xyz/quiz/service/lobby"
	"github.com/nowel-xyz/quiz/utils/web"
)

func SetupLobbyJoinRoutes(router fiber.Router) {
	router.Post("/:invitecode/join", func(c *fiber.Ctx) error {
		u := c.Locals("user")
		user, ok := u.(models.User)
		if !ok {
			log.Println("unauthorized or invalid user context")
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
		}

		lobbyInviteID := c.Params("invitecode")

		lobby, err := service_lobby.GetLobbyByInviteCode(c, lobbyInviteID, user)
		if err != nil {
			log.Println("error getting lobby by invite code:", err)
			return err // ✅ Early return to prevent nil pointer
		}

		lobbyJoined, err := service_lobby.JoinLobbyByID(c, lobby.ID, user)
		if err != nil {
			log.Println("error joining lobby:", err)
			return err // ✅ Same here
		}

		web.HubInstance.AddUserToLobby(user.ID.Hex(), lobbyJoined.ID)
		update := map[string]interface{}{
			"type":    "member-list",
			"lobbyID": lobbyJoined.ID,
			"Members": lobbyJoined.Members,
		}
		web.HubInstance.BroadcastToLobby(lobbyJoined.ID, update)
		return c.Redirect("/lobby/" + lobbyJoined.ID, fiber.StatusSeeOther)
	})
}

