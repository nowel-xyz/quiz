package service_lobby

import (
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/nowel-xyz/quiz/database"
	"github.com/nowel-xyz/quiz/database/models"
	"github.com/nowel-xyz/quiz/routers/api/lobby/utils"
)

func GetLobbyMembersById(c *fiber.Ctx, lobbyID string, user models.User) (*utils.Lobby, error) {

	// Step 2: Fetch full lobby data
	lobbyKey := "lobby:" + lobbyID
	raw, err := database.Redis.Get(c.Context(), lobbyKey).Result()
	if err != nil {
		log.Println("redis get error (lobby):", err)
		return nil, fiber.NewError(fiber.StatusNotFound, "Lobby not found")
	}

	// Step 3: Unmarshal the JSON
	var lobby utils.Lobby
	if err := json.Unmarshal([]byte(raw), &lobby); err != nil {
		log.Println("json unmarshal error:", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to parse lobby data")
	}	

	return &lobby, nil
}
