package api_auth

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	mailing "XerWoho/zep.run/internal/app/utils/mail"
	"XerWoho/zep.run/internal/database"
)

var (
	REGEXP_USERNAME = regexp.MustCompile(`^[a-zA-Z0-9_-]{3,15}$`)
)

func AuthRegisterRoute(c *gin.Context) {
	_, exists := c.Get("user")
	if exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "cannot register if already logged in",
		})
		return
	}

	type RegisterRequest struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req RegisterRequest
	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "something went wrong... parsing",
		})
		return
	}

	if !REGEXP_EMAIL.MatchString(req.Email) || !REGEXP_USERNAME.MatchString(req.Username) || utf8.RuneCountInString(req.Password) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid values...",
		})
		return
	}

	password := []byte(req.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid password...",
		})
		return
	}

	db, err := database.GetDB()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "something went wrong... db",
		})
		return
	}
	defer db.Close()

	queries := database.New(db)
	ctx := context.Background()

	createdUser, err := queries.CreateUser(ctx, database.CreateUserParams{
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		Email:     strings.ToLower(req.Email),
		Username:  req.Username,
		Password:  string(hashedPassword),
		ID:        uuid.NewString(),
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Name or E-Mail already in use.",
		})
		return
	}

	allCharacters := "0123456789"
	verificationCode := ""
	for i := 0; i < 6; i++ {
		rand.NewSource(time.Now().UnixNano())
		num := rand.Intn(len(allCharacters))
		verificationCode += string(allCharacters[num])
	}

	_, err = queries.CreateVerification(ctx, database.CreateVerificationParams{
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		ID:        uuid.NewString(),
		UserID:    createdUser.ID,
		Code:      verificationCode,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "verification could not be created...",
		})
		return
	}

	emailSuccess, err := mailing.BrevoMail(
		strings.ToLower(req.Email),
		15, // zeP register template
		map[string]interface{}{
			"USER_NAME": req.Username,
			"CODE":      verificationCode,
		},
	)

	if err != nil || !emailSuccess {
		fmt.Printf("%s", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "something went wrong...",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"error":   nil,
	})
}
