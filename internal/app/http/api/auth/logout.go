package api_auth

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"XerWoho/zep.run/internal/database"
)

func AuthLogoutRoute(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "something went wrong...",
		})
		return
	}

	userMap := user.(gin.H)
	userSessionId := userMap["session_id"].(string)

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
	err = queries.DeleteSessionById(ctx, userSessionId)
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
