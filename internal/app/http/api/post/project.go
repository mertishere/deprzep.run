package api_post

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"XerWoho/zep.run/internal/database"
)

func PostProjectRoute(c *gin.Context) {
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

	type PostProjectParams struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Docs        string `json:"docs"`
		Tags        string `json:"tags"`
	}

	var req PostProjectParams
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
	queriedProject, err := queries.CreateProject(ctx, database.CreateProjectParams{
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		ID:        uuid.NewString(),
		UserID:    userId,
		Name:      req.Name,
		Description: sql.NullString{
			String: req.Description,
			Valid:  true,
		},
		Tags: sql.NullString{
			String: req.Tags,
			Valid:  true,
		},
		Docs: sql.NullString{
			String: req.Docs,
			Valid:  true,
		},
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
		"project": queriedProject,
	})
}
