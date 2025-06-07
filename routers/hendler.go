package routers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nowel-xyz/quiz/routers/api/lobby"
    "github.com/nowel-xyz/quiz/routers/api/auth"
	"github.com/nowel-xyz/quiz/routers/web"
)

func SetupWebRoutes(app *fiber.App) {
    web.SetupQuizRoutes(app)
    web.SetupAuthRoutes(app)
    web.SetupWebSocket(app)
    web.SetupHomeRoutes(app)
    web.SetupLobbyRoutes(app)
}

func SetupAPIRoutes(v1 fiber.Router) {
    lobby.SetupLobbyRoutes(v1)


    auth.SetupAuthRoutes(v1)

}



