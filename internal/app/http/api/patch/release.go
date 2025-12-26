package api_patch

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"XerWoho/zep.run/internal/database"
)

func PatchReleaseRoute(c *gin.Context) {
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

	type PatchReleaseParams struct {
		ProjectId  string `json:"project_id"`
		Url        string `json:"url"`
		Release    string `json:"release"`
		ZigVersion string `json:"zig_version"`
		Hash       string `json:"hash"`
		RootFile   string `json:"root_file"`
	}

	var req PatchReleaseParams
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
	queriedProject, err := queries.GetProjectById(ctx, database.GetProjectByIdParams{
		ID:     req.ProjectId,
		UserID: userId,
	})

	queriedRelease, err := queries.GetReleasesByProjectIdAndRelease(ctx, database.GetReleasesByProjectIdAndReleaseParams{
		ProjectID: queriedProject.ID,
		Release:   req.Release,
	})

	_, err = queries.UpdateReleaseById(ctx, database.UpdateReleaseByIdParams{
		ID:         queriedRelease.ID,
		UserID:     userId,
		Release:    req.Release,
		Hash:       req.Hash,
		ZigVersion: req.ZigVersion,
		RootFile:   req.RootFile,
		UpdatedAt: sql.NullString{
			String: time.Now().Format("2006-01-02 15:04:05"),
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
