package api_get

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"XerWoho/zep.run/internal/database"
)

func GetReleaseRoute(c *gin.Context) {
	db, err := database.GetDB()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "something went wrong...",
		})
		return
	}
	defer db.Close()

	type GetProjectParams struct {
		Name    string `json:"name"`
		Release string `json:"release"`
	}

	var req GetProjectParams
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
	queriedProject, err := queries.GetProjectByName(ctx, req.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "something went wrong...",
		})
		return
	}

	projectId := queriedProject.ID
	queriedReleases, err := queries.GetReleasesByProjectId(ctx, projectId)
	if len(req.Release) == 0 {
		first := queriedReleases[0]
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"error":   nil,
			"project": first,
		})
		return
	}

	for i := 0; i < 10; i++ {
		r := queriedReleases[i]
		if r.Release == req.Release {
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"error":   nil,
				"project": queriedReleases[i],
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"success": false,
		"error":   "Not found.",
	})
}
