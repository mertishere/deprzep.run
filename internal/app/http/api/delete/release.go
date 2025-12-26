package api_delete

import (
	"context"
	"encoding/json"
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

	type DeleteReleaseParams struct {
		ProjectId string `json:"project_id"`
		Id        string `json:"id"`
	}

	var req DeleteReleaseParams
	err = json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "something went wrong...",
		})
		return
	}

	queries := database.New(db)
	ctx := context.Background()
	project, err := queries.GetProjectById(ctx, database.GetProjectByIdParams{
		ID:     req.ProjectId,
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
		ProjectID: req.ProjectId,
		ID:        req.Id,
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
