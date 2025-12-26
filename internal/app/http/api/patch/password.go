package api_patch

import (
	"context"
	"encoding/json"
	"net/http"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"XerWoho/zep.run/internal/database"
)

func PatchPasswordRoute(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "something went wrong...",
		})
		return
	}

	userMap := user.(gin.H)
	userId := userMap["id"].(string)

	type PasswordRequest struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	var req PasswordRequest
	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "something went wrong...",
		})
		return
	}

	if utf8.RuneCountInString(req.NewPassword) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "something went wrong...",
		})
		return
	}

	password := []byte(req.NewPassword)
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
			"error":   "something went wrong...",
		})
		return
	}
	defer db.Close()

	queries := database.New(db)
	ctx := context.Background()
	queriedUser, err := queries.GetUserById(ctx, userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "something went wrong...",
		})
		return
	}

	oldPassword := []byte(req.OldPassword)
	oldHashedPassword := []byte(queriedUser.Password)
	err = bcrypt.CompareHashAndPassword(oldHashedPassword, oldPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid password...",
		})
		return
	}

	_, err = queries.UpdateUserPassword(ctx,
		database.UpdateUserPasswordParams{
			ID: userId,

			Password: string(hashedPassword),
		},
	)

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
	})
}
