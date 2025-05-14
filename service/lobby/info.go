package service_lobby

import (
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/nowel-xyz/quiz/database"
	"github.com/nowel-xyz/quiz/database/models"
	"github.com/nowel-xyz/quiz/routers/api/lobby/utils"
)

func GetLobbyData(c *fiber.Ctx, lobbyID string, user models.User) (*utils.Lobby, error) {
	key := "lobby:" + lobbyID

	raw, err := database.Redis.Get(c.Context(), key).Result()
	if err != nil {
		log.Println("redis get error:", err)
		return nil, fiber.NewError(fiber.StatusNotFound, "Lobby not found")
	}

	var lobby utils.Lobby
	if err := json.Unmarshal([]byte(raw), &lobby); err != nil {
		log.Println("json unmarshal error:", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to parse lobby data")
	}

	if !utils.ContainsMember(lobby.Members, user.UserID) && lobby.HostID != user.UserID {
		return nil, fiber.NewError(fiber.StatusForbidden, "You're not a member of this lobby")
	}

	return &lobby, nil
}
