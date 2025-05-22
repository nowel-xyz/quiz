package web

import (
	"bytes"
	"html/template"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/nowel-xyz/quiz/database/models"
	"github.com/nowel-xyz/quiz/routers/api/lobby/utils"
	"github.com/nowel-xyz/quiz/routers/middleware"
	service_lobby "github.com/nowel-xyz/quiz/service/lobby"
)

func SetupLobbyRoutes(app *fiber.App) {
	app.Get("/lobby/:id", middleware.RequireAuth(), func(c *fiber.Ctx) error {
		tmpl := template.Must(
			template.New("lobby").
				ParseFiles("./views/lobby/index.html"),
		)
		id := c.Params("id")
		lobbyinfo, err := service_lobby.GetLobbyMembersById(c, id, models.User{})
		if err != nil {
			log.Println("error getting lobby data", err)
			return c.Status(404).SendString("Lobby not found")
		}
		var Data = struct {
			Title string
			Lobby *utils.Lobby
		}{
			Title: "Trivia Quiz",
			Lobby: lobbyinfo,
		}

		var buf bytes.Buffer
		if err := tmpl.ExecuteTemplate(&buf, "lobby", Data); err != nil {
			log.Println("Template exec error:", err)
			return c.Status(500).SendString("Template execution error: " + err.Error())
		}
		return c.Type("html").Send(buf.Bytes())
	})

	app.Get("/lobby/join", func(c *fiber.Ctx) error {
		tmpl := template.Must(
			template.New("lobby").
				ParseFiles("./views/lobby/join.html"),
		)

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

	app.Post("/lobby/host", func(c *fiber.Ctx) error {
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
