package auth

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"

	"github.com/nowel-xyz/quiz/database"
	"github.com/nowel-xyz/quiz/database/models"
	"github.com/nowel-xyz/quiz/routers/api/auth/utils"
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

func SetupLoginRoutes(router fiber.Router) {
	users := database.Database.Collection("users")

	router.Post("/login", func(c *fiber.Ctx) error {
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

		_, err = users.UpdateOne(context.TODO(), bson.M{"_id": user.ID}, bson.M{"$set": bson.M{"ips": user.IPs}})
		if err != nil {
			return c.Status(500).SendString("db error")
		}

		token := user.Cookie
		if token == "" {
			j := jwt.NewWithClaims(jwt.SigningMethodHS256, 
				jwt.MapClaims{
					"id": user.ID, 
					"email": user.Email, 
					"exp": time.Now().Add(7 * 24 * time.Hour).Unix(),
				},
			)
			
			token, _ = j.SignedString(utils.GetJWTKey())
			_, _ = users.UpdateOne(context.TODO(), bson.M{"_id": user.ID}, bson.M{"$set": bson.M{"cookie": token}})
		}

		c.Cookie(&fiber.Cookie{
			Name:     "sessionToken",
			Value:    token,
			HTTPOnly: true,
			Secure:   false,
			Expires:  time.Now().Add(7 * 24 * time.Hour),
		})

		go sendMail(user.Email, subject, text)

		return c.Redirect("/", fiber.StatusSeeOther)
	})
}
