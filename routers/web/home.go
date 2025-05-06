package web

import (
	"bytes"
	"context"
	"html/template"
	"log"


	"github.com/gofiber/fiber/v2"
	"github.com/nowel-xyz/quiz/database"
	"github.com/nowel-xyz/quiz/database/models"
	"go.mongodb.org/mongo-driver/bson"
)

func SetupHomeRoutes(app *fiber.App) {
	users := database.Database.Collection("users")

	tmpl := template.Must(
		template.New("home").
			ParseFiles(
				"./views/header.html",
				"./views/index.html",
			),
	)

	app.Get("/", func(c *fiber.Ctx) error {
		type Data struct {
			Title     string
			Username  string
			Roles     []string
			LoggedIn  bool
		}

		// Initializing default data
		data := Data{
			Title:    "Home",
			LoggedIn: false, // Default to not logged in
		}

		// Get cookie from browser
		tokenStr := c.Cookies("sessionToken")
		log.Println("Token from cookie:", tokenStr)

		var user models.User
        err := users.FindOne(context.TODO(), bson.M{"cookie": tokenStr}).Decode(&user)
        log.Println("User from DB:", user)
        if err == nil {
            data.Username = user.Username
            data.Roles = user.Roles
            data.LoggedIn = true
        } else {
            log.Println("User lookup error:", err)
        }

		// Render HTML template
		var buf bytes.Buffer
		if err := tmpl.ExecuteTemplate(&buf, "home", data); err != nil {
			log.Println("Template exec error:", err)
			return c.Status(500).SendString("Template execution error: " + err.Error())
		}

		return c.Type("html").Send(buf.Bytes())
	})
}
