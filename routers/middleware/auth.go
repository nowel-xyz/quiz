package middleware

import (
	"context"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nowel-xyz/quiz/database"
	"github.com/nowel-xyz/quiz/database/models"
	"go.mongodb.org/mongo-driver/bson"
)


func RequireAuth() fiber.Handler {
    jwtKey := []byte(os.Getenv("SECRET_KEY"))
    usersColl := database.Database.Collection("users")

    return func(c *fiber.Ctx) error {
        // 1) Get token from cookie
        tok := c.Cookies("sessionToken")
        if tok == "" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
        }

        // 2) Parse & validate JWT
        token, err := jwt.Parse(tok, func(t *jwt.Token) (interface{}, error) {
            if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fiber.ErrUnauthorized
            }
            return jwtKey, nil
        })

        if err != nil || !token.Valid {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "invalid token"})
        }
		

        // 3) Extract user ID claim
        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "invalid claims"})
        }

                
        userID, ok := claims["id"].(string)
        if !ok {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "invalid claims data"})
        }

        // 4) Fetch the user document once
        var user models.User
        err = usersColl.FindOne(context.TODO(), bson.M{"id": userID}).Decode(&user)
        if err != nil {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "user not found"})
        }

        // 5) Store in locals for downstream handlers
        c.Locals("user", user)

        return c.Next()
    }
}
