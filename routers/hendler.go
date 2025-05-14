package routers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nowel-xyz/quiz/routers/api"
	"github.com/nowel-xyz/quiz/routers/api/lobby"
	"github.com/nowel-xyz/quiz/routers/web"
	"github.com/nowel-xyz/quiz/routers/middleware"
)

func SetupWebRoutes(app *fiber.App) {
    web.SetupQuizRoutes(app)
    web.SetupAuthRoutes(app)
    web.SetupWebSocket(app)
    web.SetupHomeRoutes(app)
    web.SetupLobbyRoutes(app)
}

func SetupAPIRoutes(v1 fiber.Router) {
    api.SetupAuthRoutes(v1)

	lobbyRouter := v1.Group("/lobby", middleware.RequireAuth())
    lobby.SetupLobbyRoutes(lobbyRouter)

}



