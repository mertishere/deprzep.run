package api_auth

import (
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	security "XerWoho/zep.run/internal/app/security"
	"XerWoho/zep.run/internal/database"
)

var (
	REGEXP_EMAIL = regexp.MustCompile(`^[a-zA-Z0-9_+&*-]+(?:\.[a-zA-Z0-9_+&*-]+)*@(?:[a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$`)
)

func AuthLoginRoute(c *gin.Context) {
	_, exists := c.Get("user")
	if exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "already logged in",
		})
		return
	}

	type LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LoginRequest
	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "something went wrong...",
		})
		return
	}

	if !REGEXP_EMAIL.MatchString(req.Email) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid values...",
		})
		return
	}

	db, err := database.GetDB()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "something went wrong...",
		})
		return
	}
	defer db.Close()

	queries := database.New(db)
	ctx := context.Background()
	user, err := queries.GetUserByEmail(ctx, strings.ToLower(req.Email))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "not found.",
		})
		return
	}

	if !user.IsVerified {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"success": false,
			"error":   "email verification required.",
		})
		return
	}

	password := []byte(req.Password)
	hashedPassword := []byte(user.Password)
	err = bcrypt.CompareHashAndPassword(hashedPassword, password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid password...",
		})
		return
	}

	allCharacters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890-_"
	sessionSecret := ""
	for i := 0; i < 64; i++ {
		rand.NewSource(time.Now().UnixNano())
		num := rand.Intn(64)
		sessionSecret += string(allCharacters[num])
	}

	session, err := queries.CreateSession(ctx, database.CreateSessionParams{
		ID:        uuid.NewString(),
		UserID:    user.ID,
		Secret:    sessionSecret,
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "something went wrong...",
		})
		return
	}

	jwt, err := security.JWTEncode(security.JWTStruct{
		Name:            user.Username,
		CreatedAt:       user.CreatedAt,
		PublicUserId:    user.ID,
		PublicSessionId: session.ID,
		SessionSecret:   session.Secret,
		Email:           user.Email,
	}, os.Getenv("SECRET_JWT_KEY"), nil)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "something went wrong...",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"error":   nil,
		"jwt":     jwt,
	})
}
