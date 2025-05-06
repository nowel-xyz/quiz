package lobby

import (
	"encoding/json"
	"log"

	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/nowel-xyz/quiz/database"
	"github.com/nowel-xyz/quiz/database/models"
	"github.com/nowel-xyz/quiz/routers/api/lobby/utils"
)

func SetupLobbyCreateRoutes(router fiber.Router) {

	type Settings struct {
		MemberLimit int `json:"member_limit"`
	}

	type LobbyInvite struct {
		Code string `json:"code"`
	}

	type Lobby struct {
		ID        string      `json:"lobby_id"`
		Invite    LobbyInvite `json:"lobby_invite"`
		HostID    string      `json:"host_id"`
		QuizID    string      `json:"quiz_id"`
		Members   []string    `json:"members"`
		Settings  Settings    `json:"settings"`
		CreatedAt string      `json:"created_at"`
		UpdatedAt string      `json:"updated_at"`
	}

	router.Post("/create", func(c *fiber.Ctx) error {
		user := c.Locals("user").(models.User)
		type req struct {
			QuizID   string   `json:"quiz_id"`
			Settings Settings `json:"settings"`
		}

		var body req
		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).SendString("bad request")
		}

		code, err := utils.GenerateCode(10)
		if err != nil {
			log.Fatal(err)
		}

		newLobby := Lobby{
			ID: uuid.NewString(),
			Invite: LobbyInvite{
				Code: code,
			},
			HostID:    user.UserID,
			QuizID:    body.QuizID,
			Members:   []string{user.UserID},
			Settings:  body.Settings,
			CreatedAt: time.Now().Format(time.RFC3339),
			UpdatedAt: time.Now().Format(time.RFC3339),
		}  

		key := "lobby:" + newLobby.Invite.Code

		// 2) JSON‑marshal the lobby struct
		blob, err := json.Marshal(newLobby)
		if err != nil {
			log.Println("json marshal error:", err)
			return c.Status(fiber.StatusInternalServerError).SendString("internal error")
		}
	
		// 3) Store it under that key
		if err := database.Redis.Set(c.Context(), key, blob, 24*time.Hour).Err(); err != nil {
			log.Println("redis set error:", err)
			return c.Status(fiber.StatusInternalServerError).SendString("redis error")
		}
	
		// 4) Read it back (for testing/logging)
		raw, err := database.Redis.Get(c.Context(), key).Result()
		if err != nil {
			log.Println("redis get error:", err)
		} else {
			var stored Lobby
			if err := json.Unmarshal([]byte(raw), &stored); err != nil {
				log.Println("json unmarshal error:", err)
			} else {
				log.Println("Lobby from Redis:", stored)
			}
		}
	
		return c.Status(fiber.StatusCreated).JSON(newLobby)
	})
}
