package web

import (
	"bytes"
	"html/template"
	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(app *fiber.App) {
	tmpl := template.Must(
		template.New("auth").
			ParseFiles(
				"./views/header.html",
				"./views/auth.html",
			),
	)

	// Render the login page
	app.Get("/login", func(c *fiber.Ctx) error {
		var buf bytes.Buffer
		if err := tmpl.ExecuteTemplate(&buf, "login", nil); err != nil {
			return c.Status(500).SendString("Template execution error: " + err.Error())
		}
		return c.Type("html").Send(buf.Bytes())
	})

	// Render the register page
	app.Get("/register", func(c *fiber.Ctx) error {
		var buf bytes.Buffer
		if err := tmpl.ExecuteTemplate(&buf, "register", nil); err != nil {
			return c.Status(500).SendString("Template execution error: " + err.Error())
		}
		return c.Type("html").Send(buf.Bytes())
	})
}
