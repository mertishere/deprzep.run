package api_get

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"XerWoho/zep.run/internal/database"
)

func GetProjectRoute(c *gin.Context) {
	db, err := database.GetDB()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "something went wrong...",
		})
		return
	}
	defer db.Close()

	name := c.Query("name")

	queries := database.New(db)
	ctx := context.Background()
	queriedProject, err := queries.GetProjectByName(ctx, name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "something went wrong...",
		})
		return
	}

	type Project struct {
		ID          string
		UserID      string
		Name        string
		Description string
		Docs        string
		Tags        string
		CreatedAt   string
	}

	r := Project{
		queriedProject.ID,
		queriedProject.UserID,
		queriedProject.Name,
		queriedProject.Description.String,
		queriedProject.Docs.String,
		queriedProject.Tags.String,
		queriedProject.CreatedAt,
	}

	queriedReleases, err := queries.GetReleasesByProjectId(ctx, r.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "something went wrong...",
		})
		return
	}

	type Release struct {
		ID         string
		UserID     string
		ProjectID  string
		Url        string
		Release    string
		ZigVersion string
		Hash       string
		RootFile   string
		CreatedAt  string
		UpdatedAt  string
	}

	parsedReleases := []Release{}
	for i := 0; i < len(queriedReleases); i++ {
		release := queriedReleases[i]
		v := Release{
			release.ID,
			release.UserID,
			release.ProjectID,
			release.Url,
			release.Release,
			release.ZigVersion,
			release.Hash,
			release.RootFile,
			release.CreatedAt,
			release.UpdatedAt.String,
		}
		parsedReleases = append(parsedReleases, v)
	}

	jProject, err := json.Marshal(r)
	jReleased, err := json.Marshal(parsedReleases)

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"error":    nil,
		"project":  jProject,
		"releases": jReleased,
	})
}

func GetProjectsRoute(c *gin.Context) {
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
	queriedProjects, err := queries.GetProjectsByUserId(ctx, userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "something went wrong...",
		})
		return
	}

	type Project struct {
		ID          string
		UserID      string
		Name        string
		Description string
		Docs        string
		Tags        string
		CreatedAt   string
	}

	parsedProject := []Project{}
	for i := 0; i < len(queriedProjects); i++ {
		project := queriedProjects[i]
		r := Project{
			project.ID,
			project.UserID,
			project.Name,
			project.Description.String,
			project.Docs.String,
			project.Tags.String,
			project.CreatedAt,
		}
		parsedProject = append(parsedProject, r)
	}

	jProjects, err := json.Marshal(parsedProject)
	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"error":    nil,
		"projects": jProjects,
	})
}
