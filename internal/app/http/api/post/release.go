package api_post

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"XerWoho/zep.run/internal/buckets"
	"XerWoho/zep.run/internal/database"
)

func PostReleaseRoute(c *gin.Context) {
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

	type PostReleaseParams struct {
		ProjectId  string `json:"project_id"`
		Hash       string `json:"hash"`
		Release    string `json:"release"`
		ZigVersion string `json:"zig_version"`
		RootFile   string `json:"root_file"`
	}

	r_id := c.PostForm("project_id")
	r_hash := c.PostForm("hash")
	r_release := c.PostForm("release")
	r_zig_version := c.PostForm("zig_version")
	r_root_file := c.PostForm("root_file")
	req := PostReleaseParams{
		r_id,
		r_hash,
		r_release,
		r_zig_version,
		r_root_file,
	}

	formFile, err := c.FormFile("package")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "something went wrong...",
		})
		return
	}

	if formFile.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "something went wrong...",
		})
		return
	}

	openedFile, err := formFile.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "something went wrong...",
		})
		return
	}
	defer openedFile.Close()
	mimeTypeBuffer := make([]byte, 512)
	_, err = openedFile.Read(mimeTypeBuffer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "something went wrong...",
		})
		return
	}
	mimeType := http.DetectContentType(mimeTypeBuffer)
	validFileExtensions := []string{
		"application/zstd",
		"application/octet-stream",
	}

	if !slices.Contains(validFileExtensions, mimeType) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "something went wrong...",
		})
		return
	}

	splitFileName := strings.Split(formFile.Filename, ".")
	fileExtension := strings.Join(splitFileName[1:], ".")
	if fileExtension != "tar.zstd" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "something went wrong...",
		})
		return
	}

	// create temporary file
	f, err := os.CreateTemp("", fmt.Sprintf("sample-*.%s", fileExtension))
	defer os.Remove(f.Name())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "something went wrong...",
		})
		return
	}

	// save it to temp path
	err = c.SaveUploadedFile(formFile, f.Name())
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
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "something went wrong...",
		})
		return
	}

	uploadNewFile := fmt.Sprintf("projects/%s/%s-%s-%s.%s", userId, queriedProject.Name, req.Release, req.Hash, fileExtension)
	url, err := buckets.UploadFile(
		f.Name(),
		uploadNewFile,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "something went wrong...",
		})
		return
	}

	queriedRelease, err := queries.GetReleasesByProjectIdAndRelease(ctx, database.GetReleasesByProjectIdAndReleaseParams{
		ProjectID: req.ProjectId,
		Release:   req.Release,
	})

	if err == nil {
		if queriedRelease.UserID == userId {
			err = queries.DeleteRelease(ctx, database.DeleteReleaseParams{
				ProjectID: req.ProjectId,
				UserID:    userId,
				Release:   req.Release,
			})
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"error":   "something went wrong...",
				})
				return
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "something went wrong...",
			})
			return
		}
	}

	createRelease, err := queries.CreateRelease(ctx, database.CreateReleaseParams{
		CreatedAt:  time.Now().Format("2006-01-02 15:04:05"),
		ID:         uuid.NewString(),
		UserID:     userId,
		ProjectID:  req.ProjectId,
		Hash:       req.Hash,
		Release:    req.Release,
		ZigVersion: req.ZigVersion,
		RootFile:   req.RootFile,
		Url:        url,
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
		"release": createRelease,
	})
}
