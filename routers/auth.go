package routers

import (
	"bytes"
	"context"
	"html/template"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"

	"github.com/nowel-xyz/quiz/database"
	"github.com/nowel-xyz/quiz/database/models"
)

func sendMail(to, subject, body string) error {
	mailUser := os.Getenv("MAIL_USER")
	mailPass := os.Getenv("MAIL_PASSWORD")

	m := gomail.NewMessage()
	m.SetHeader("From", mailUser)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer("smtp.gmail.com", 587, mailUser, mailPass)
	return d.DialAndSend(m)
}

func SetupAuthRoutes(app *fiber.App) {
	users := database.Database.Collection("users")
	jwtKey := []byte(os.Getenv("SECRET_KEY"))

	// Setup templates for login and register pages
	tmpl := template.Must(
		template.New("auth").
			ParseFiles(
				"./views/header.html",
				"./views/auth.html", // The auth template
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

	// Handle registration logic
	app.Post("/register", func(c *fiber.Ctx) error {
		type req struct {
			Username string `json:"username"`
			Password string `json:"password"`
			Email    string `json:"email"`
		}
		var body req
		if err := c.BodyParser(&body); err != nil || body.Email == "" || body.Password == "" || body.Username == "" {
			return c.Status(400).SendString("bad request")
		}

		body.Email = strings.ToLower(body.Email)
		body.Username = strings.ToLower(body.Username)

		count, err := users.CountDocuments(context.TODO(), bson.M{"email": body.Email})
		if err != nil || count > 0 {
			return c.Status(400).SendString("email already taken")
		}

		id := uuid.NewString()
		salt, _ := bcrypt.GenerateFromPassword([]byte(time.Now().String()), 12)
		hashed, _ := bcrypt.GenerateFromPassword([]byte(body.Password), 12)

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":    id,
			"email": body.Email,
		})
		signedToken, _ := token.SignedString(jwtKey)

		ip := c.IP()
		user := models.User{
			ID:       primitive.NewObjectID(),
			UserID:   id,
			Username: body.Username,
			Password: string(hashed),
			Salt:     string(salt),
			Email:    body.Email,
			Cookie:   signedToken,
			IPs: []models.IPEntry{{
				IP:         ip,
				LoginTimes: 0,
				LastLogin:  time.Now(),
			}},
		}

		if _, err := users.InsertOne(context.TODO(), user); err != nil {
			return c.Status(500).SendString("failed to save")
		}

		go sendMail(body.Email, "Registration", "Welcome to our platform, "+body.Username)

		return c.JSON(fiber.Map{"message": "registered"})
	})

	// Handle login logic
	app.Post("/login", func(c *fiber.Ctx) error {
		type req struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		var body req
		if err := c.BodyParser(&body); err != nil || body.Email == "" || body.Password == "" {
			return c.Status(400).SendString("bad request")
		}

		body.Email = strings.ToLower(body.Email)

		var user models.User
		err := users.FindOne(context.TODO(), bson.M{"email": body.Email}).Decode(&user)
		if err != nil {
			return c.Status(401).SendString("unauthorized")
		}

		if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)) != nil {
			return c.Status(401).SendString("unauthorized")
		}

		ip := c.IP()
		found := false
		subject := ""
		text := ""

		for i, entry := range user.IPs {
			if entry.IP == ip {
				user.IPs[i].LoginTimes += 1
				user.IPs[i].LastLogin = time.Now()
				found = true
				subject = "New login from recognized IP"
				text = "New login from known IP address."
				break
			}
		}

		if !found {
			user.IPs = append(user.IPs, models.IPEntry{IP: ip, LoginTimes: 1, LastLogin: time.Now()})
			subject = "New login from unrecognized IP"
			text = "A new login attempt was made from a new IP: " + ip
		}

		// update user record
		_, err = users.UpdateOne(context.TODO(), bson.M{"_id": user.ID}, bson.M{"$set": bson.M{"ips": user.IPs}})
		if err != nil {
			return c.Status(500).SendString("db error")
		}

		token := user.Cookie
		if token == "" {
			j := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"id": user.UserID, "email": user.Email})
			token, _ = j.SignedString(jwtKey)
			_, _ = users.UpdateOne(context.TODO(), bson.M{"_id": user.ID}, bson.M{"$set": bson.M{"cookie": token}})
		}

		// set token in cookie
		c.Cookie(&fiber.Cookie{
			Name:     "sessionToken",
			Value:    token,
			HTTPOnly: true,
			Secure:   true,
			Expires:  time.Now().Add(7 * 24 * time.Hour),
		})

		go sendMail(user.Email, subject, text)

		return c.JSON(fiber.Map{
			"message": "login success",
			"user": fiber.Map{
				"username": user.Username,
				"roles":    user.Roles,
			},
		})
	})

	// Logout
	app.Get("/logout", func(c *fiber.Ctx) error {
		c.ClearCookie("sessionToken")
		return c.Redirect(os.Getenv("FRONTEND") + "/")
	})
}
