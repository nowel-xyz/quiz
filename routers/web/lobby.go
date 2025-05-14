package web

import (
	"bytes"
	"html/template"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/nowel-xyz/quiz/database/models"
	"github.com/nowel-xyz/quiz/routers/middleware"
	service_lobby "github.com/nowel-xyz/quiz/service/lobby"
)

func SetupLobbyRoutes(app *fiber.App) {
	app.Get("/lobby", func(c *fiber.Ctx) error {
	
		tmpl := template.Must(
			template.New("lobby").
				ParseFiles("./views/lobby/lobby.html"),
		)

		var Data = struct {
			Title string
			
		}{
			Title: "Trivia Quiz",
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

	app.Get("/lobby/:id", middleware.RequireAuth(), func(c *fiber.Ctx) error {
		
		tmpl := template.Must(
			template.New("lobby_info").
				ParseFiles("./views/lobby/info.html"),
		)

		id := c.Params("id")
	
		userVal := c.Locals("user")
		if userVal == nil {
			return c.Status(fiber.StatusUnauthorized).SendString("Not logged in")
		}
	
		user, ok := userVal.(models.User)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).SendString("User context invalid")
		}
	
		lobbyData, err := service_lobby.GetLobbyData(c, id, user)
		if err != nil {
			log.Println("error getting lobby:", err)
			return c.Status(404).SendString("Lobby not found")
		}
	
		var Data = struct {
			Title string
			Lobby any
		}{
			Title: "Trivia Quiz",
			Lobby: lobbyData,
		}
	
		var buf bytes.Buffer
		if err := tmpl.ExecuteTemplate(&buf, "lobby_info", Data); err != nil {
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
