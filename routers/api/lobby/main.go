package lobby

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nowel-xyz/quiz/routers/middleware"

)

func SetupLobbyRoutes(router fiber.Router) {
	lobbyRouter := router.Group("/lobby", middleware.RequireAuth())
	SetupLobbyCreateRoutes(lobbyRouter)
	SetupLobbyInfoRoutes(lobbyRouter)
	SetupLobbyJoinRoutes(lobbyRouter)
}
