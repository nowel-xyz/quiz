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

	router.Post("/create", func(c *fiber.Ctx) error {
		user := c.Locals("user").(models.User)
		type req struct {
			QuizID   string   `json:"quiz_id"`
			Settings utils.Settings `json:"settings"`
		}

		body := new(req)
		if c.Request().Body() != nil && len(c.Request().Body()) > 0 {
			if err := c.BodyParser(body); err != nil {
				return c.Status(400).SendString("bad request")
			}
}
		code, err := utils.GenerateCode(10)
		if err != nil {
			log.Fatal(err)
		}
		lobbyId := uuid.NewString()

		newLobby := utils.Lobby{
			ID: lobbyId,
			Invite: utils.LobbyInvite{
				Code: code,
			},
			HostID:    user.UserID,
			QuizID:    body.QuizID,
			Settings:  body.Settings,
			CreatedAt: time.Now().Format(time.RFC3339),
			UpdatedAt: time.Now().Format(time.RFC3339),
		}  

		lobbyKey := "lobby:" + lobbyId
		inviteKey := "invite:" + code

		// Marshal lobby data
		blob, err := json.Marshal(newLobby)
		if err != nil {
			log.Println("json marshal error:", err)
			return c.Status(fiber.StatusInternalServerError).SendString("internal error")
		}

		// Save lobby to Redis
		if err := database.Redis.Set(c.Context(), lobbyKey, blob, 24*time.Hour).Err(); err != nil {
			log.Println("redis set error (lobby):", err)
			return c.Status(fiber.StatusInternalServerError).SendString("redis save lobby error")
		}

		// Save invite code → lobby ID mapping
		if err := database.Redis.Set(c.Context(), inviteKey, lobbyId, 24*time.Hour).Err(); err != nil {
			log.Println("redis set error (invite code):", err)
			return c.Status(fiber.StatusInternalServerError).SendString("redis save invite code error")
		}

		log.Printf("Lobby created: lobbyId=%s inviteCode=%s", lobbyId, code)

		return c.Redirect("/lobby/" + lobbyId) 
	})
}
