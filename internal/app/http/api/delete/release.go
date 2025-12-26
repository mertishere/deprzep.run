package api_delete

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"XerWoho/zep.run/internal/database"
)

func DeleteReleaseRoute(c *gin.Context) {
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

	releaseId := c.Query("id")
	projectId := c.Query("project_id")

	queries := database.New(db)
	ctx := context.Background()
	project, err := queries.GetProjectById(ctx, database.GetProjectByIdParams{
		ID:     projectId,
		UserID: userId,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "something went wrong...",
		})
		return
	}

	if project.UserID != userId {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "not authorized...",
		})
		return
	}

	err = queries.DeleteReleaseById(ctx, database.DeleteReleaseByIdParams{
		UserID:    userId,
		ProjectID: projectId,
		ID:        releaseId,
	})
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
