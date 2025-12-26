package middleware

import (
	"context"
	"os"
	"strings"

	security "XerWoho/zep.run/internal/app/security"
	"XerWoho/zep.run/internal/database"

	"github.com/gin-gonic/gin"
)

func ValidateSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		split_auth := strings.Split(auth, " ")
		if len(split_auth) < 2 {
			c.Next()
			return
		}

		key := split_auth[0]
		if key != "Bearer" {
			c.Next()
			return
		}

		jwt := split_auth[1]

		decodedJwtCookie, err := security.JWTDecode(jwt, os.Getenv("SECRET_JWT_KEY"))
		if err != nil {
			c.Next()
			return
		}

		c.Set("user", gin.H{
			"id":         decodedJwtCookie.PublicUserId,
			"createdAt":  decodedJwtCookie.CreatedAt,
			"session_id": decodedJwtCookie.PublicSessionId,
			"email":      decodedJwtCookie.Email,
			"name":       decodedJwtCookie.Name,
		})

		c.Next()

		db, err := database.GetDB()
		if err != nil {
			c.Next()
			return
		}
		defer db.Close()

		queries := database.New(db)
		ctx := context.Background()

		session, err := queries.GetSessionById(ctx, database.GetSessionByIdParams{
			ID:     decodedJwtCookie.PublicSessionId,
			UserID: decodedJwtCookie.PublicUserId,
		})
		if err != nil {
			c.Next()
			return
		}

		if session.Secret != decodedJwtCookie.SessionSecret {
			c.Next()
			return
		}

		if session.UserID != decodedJwtCookie.PublicUserId {
			c.Next()
			return
		}
		c.Set("user", gin.H{
			"id":         decodedJwtCookie.PublicUserId,
			"createdAt":  decodedJwtCookie.CreatedAt,
			"session_id": decodedJwtCookie.PublicSessionId,
			"email":      decodedJwtCookie.Email,
			"name":       decodedJwtCookie.Name,
		})
		c.Next()
	}
}
