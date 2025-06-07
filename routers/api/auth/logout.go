package auth

import (
	"os"
	"github.com/gofiber/fiber/v2"
)



func SetupLogoutRoutes(router fiber.Router) {
	router.Get("/logout", func(c *fiber.Ctx) error {
		c.ClearCookie("sessionToken")
		return c.Redirect(os.Getenv("FRONTEND") + "/")
	})
}
