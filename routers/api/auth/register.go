package auth

import (
	"context"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"

	"github.com/nowel-xyz/quiz/database"
	"github.com/nowel-xyz/quiz/database/models"
	"github.com/nowel-xyz/quiz/routers/api/auth/utils"

)

func SetupRegisterRoutes(router fiber.Router) {
	users := database.Database.Collection("users")

	router.Post("/register", func(c *fiber.Ctx) error {
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
		signedToken, _ := token.SignedString(utils.GetJWTKey())

		ip := c.IP()
		user := models.User{
			ObjID:       primitive.NewObjectID(),
			ID:   id,
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

		return c.Redirect("/login", fiber.StatusSeeOther)
	})
}
