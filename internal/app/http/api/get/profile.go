package api_get

import (
	"context"
	"encoding/json"
	"net/http"

	"XerWoho/zep.run/internal/database"

	"github.com/gin-gonic/gin"
)

func GetWhoamiRoute(c *gin.Context) {
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

	type User struct {
		Id        string
		Username  string
		Email     string
		CreatedAt string
	}

	qUser := User{
		queriedUser.ID,
		queriedUser.Username,
		queriedUser.Email,
		queriedUser.CreatedAt,
	}
	jUser, err := json.Marshal(qUser)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"error":   nil,
		"user":    jUser,
	})
}

func GetUserRoute(c *gin.Context) {
	userName := c.Query("name")
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

	_, err = queries.GetUserByName(ctx, userName)
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

func GetEmailInUseRoute(c *gin.Context) {
	userEmail := c.Query("email")
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

	_, err = queries.GetUserByEmail(ctx, userEmail)
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
