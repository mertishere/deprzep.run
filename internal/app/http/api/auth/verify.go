package api_auth

import (
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	security "XerWoho/zep.run/internal/app/security"
	"XerWoho/zep.run/internal/database"
)

func AuthVerifyRoute(c *gin.Context) {
	_, exists := c.Get("user")
	if exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "dude you are already logged in, tf",
		})
		return
	}

	type VerifyRequest struct {
		Code  string `json:"code"`
		Email string `json:"email"`
	}

	var req VerifyRequest
	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "something went wrong... parsing",
		})
		return
	}

	verificationCode := req.Code
	userEmail := req.Email

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

	queriedUser, err := queries.GetUserByEmail(ctx, userEmail)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "something went wrong... found",
		})
		return
	}

	queriedUserVerify, err := queries.GetVerificationByCode(ctx, database.GetVerificationByCodeParams{
		UserID: queriedUser.ID,
		Code:   verificationCode,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "something went wrong... c",
		})
		return
	}

	if queriedUserVerify.UserID != queriedUser.ID {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "something went wrong... m",
		})
		return
	}
	err = queries.DeleteVerification(
		ctx,
		database.DeleteVerificationParams{
			UserID: queriedUser.ID,
			Code:   verificationCode,
		},
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "something went wrong...",
		})
		return
	}

	queries.UpdateUserVerify(ctx, database.UpdateUserVerifyParams{
		ID:         queriedUser.ID,
		IsVerified: true,
	})

	allCharacters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890-_"
	sessionSecret := ""
	for i := 0; i < 64; i++ {
		rand.NewSource(time.Now().UnixNano())
		num := rand.Intn(64)
		sessionSecret += string(allCharacters[num])
	}

	session, err := queries.CreateSession(ctx, database.CreateSessionParams{
		ID:        uuid.NewString(),
		UserID:    queriedUser.ID,
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
		Name:            queriedUser.Username,
		CreatedAt:       queriedUser.CreatedAt,
		PublicUserId:    queriedUser.ID,
		PublicSessionId: session.ID,
		SessionSecret:   session.Secret,
		Email:           queriedUser.Email,
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
