package web

import (
	"bytes"
	"html/template"
	"log"

	"github.com/gofiber/fiber/v2"

)

func SetupLobbyRoutes(app *fiber.App) {
	tmpl := template.Must(
		template.New("lobby").
			ParseFiles(
				"./views/lobby.html",
			),
	)


	app.Get("/lobby/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		var Data = struct {
			Title string
			Id    string
		}{
			Title: "Trivia Quiz",
			Id:    id,
		}

		var buf bytes.Buffer
		if err := tmpl.ExecuteTemplate(&buf, "lobby", Data); err != nil {
			log.Println("Template exec error:", err)
			return c.Status(500).SendString("Template execution error: " + err.Error())
		}
		return c.Type("html").Send(buf.Bytes())
	})


	app.Post("/lobby/create", func(c *fiber.Ctx) error {
		type req struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		var body req
		if err := c.BodyParser(&body); err != nil || body.Email == "" || body.Password == "" {
			return c.Status(400).SendString("bad request")
		}



		return c.JSON(fiber.Map{
			"message": "login success",
			
		})
	

	})

}
