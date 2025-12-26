package api_delete

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"XerWoho/zep.run/internal/buckets"
	"XerWoho/zep.run/internal/database"
)

func DeleteAccountRoute(c *gin.Context) {
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

	err = queries.DeleteSessionByUserId(ctx, userId)
	if err != nil {
	}

	_, err = buckets.DeleteFolder(fmt.Sprintf("releases/%s", userId))
	if err != nil {
	}
	err = queries.DeleteUser(ctx, userId)
	if err != nil {
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"error":   nil,
	})
}
