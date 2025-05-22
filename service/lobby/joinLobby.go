package service_lobby

import (
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/nowel-xyz/quiz/database"
	"github.com/nowel-xyz/quiz/database/models"
	"github.com/nowel-xyz/quiz/routers/api/lobby/utils"
)

func JoinLobbyByID(c *fiber.Ctx, lobbyID string, user models.User) (*utils.Lobby, error) {
	lobbyKey := "lobby:" + lobbyID

	// Step 1: Fetch lobby
	raw, err := database.Redis.Get(c.Context(), lobbyKey).Result()
	if err != nil {
		log.Println("redis get error (lobby):", err)
		return nil, fiber.NewError(fiber.StatusNotFound, "Lobby not found")
	}

	var lobby utils.Lobby
	if err := json.Unmarshal([]byte(raw), &lobby); err != nil {
		log.Println("json unmarshal error:", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to parse lobby data")
	}

	// Step 2: Add user if not already a member
	memberIDs := make([]string, len(lobby.Members))
	for i, m := range lobby.Members {
		memberIDs[i] = m.UserID
	}
	if !utils.ContainsMember(memberIDs, user.UserID) {
		lobby.Members = append(lobby.Members, models.LobbyUser{
			UserID:   user.UserID,
			Username: user.Username,
			Email:    user.Email,
			Roles:    user.Roles,
		})


		// Step 3: Marshal and store updated lobby
		updated, err := json.Marshal(lobby)
		if err != nil {
			log.Println("json marshal error:", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to save lobby")
		}

		if err := database.Redis.Set(c.Context(), lobbyKey, updated, 0).Err(); err != nil {
			log.Println("redis set error:", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to store lobby")
		}
	} else {
		return nil, fiber.NewError(fiber.StatusConflict, "Already a member of this lobby")
	}

	return &lobby, nil
}
