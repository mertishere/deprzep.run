package server

import (
	"fmt"
	"os"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"

	api_auth "XerWoho/zep.run/internal/app/http/api/auth"
	api_delete "XerWoho/zep.run/internal/app/http/api/delete"
	api_get "XerWoho/zep.run/internal/app/http/api/get"
	api_patch "XerWoho/zep.run/internal/app/http/api/patch"
	api_post "XerWoho/zep.run/internal/app/http/api/post"

	middleware "XerWoho/zep.run/internal/app/http/middleware"
	security "XerWoho/zep.run/internal/app/security"
)

func RunServer() error {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Use(gin.Recovery())
	r.StaticFile("/download.json", "./static/download.json")
	r.Static("/releases", "./static/releases")
	r.Static("/installer", "./static/installer")
	r.Static("/packages", "./static/packages")

	r.Use(gzip.Gzip(gzip.DefaultCompression))

	r.Use(middleware.ValidateSession())
	r.Use(middleware.LoggingHandler())

	r.Use(security.CORSHandler())
	r.Use(security.RateLimiter(security.RequestAmountTypes{
		AuthedUsers: 20,
		GuestUsers:  20,
	}, 10))

	api := r.Group("/api")

	api.GET("/health", api_get.GetHealth)
	api.GET("/whoami", api_get.GetWhoamiRoute)
	api.GET("/get/name", api_get.GetUserRoute)
	api.GET("/get/email", api_get.GetEmailInUseRoute)
	api.GET("/get/projects", api_get.GetProjectsRoute)
	api.GET("/get/project", api_get.GetProjectRoute)
	api.GET("/get/project/release", api_get.GetProjectRoute)

	api.PATCH("/patch/password", security.RateLimiter(security.RequestAmountTypes{
		AuthedUsers: 3,
		GuestUsers:  1,
	}, 60*1), api_patch.PatchPasswordRoute)

	api.POST("/post/project", api_post.PostProjectRoute)
	api.POST("/post/release", api_post.PostReleaseRoute)

	api.POST("/patch/project", api_patch.PatchProjectRoute)
	api.POST("/patch/release", api_patch.PatchReleaseRoute)

	api.DELETE("/delete/project", api_delete.DeleteProjectRoute)
	api.DELETE("/delete/release", api_delete.DeleteReleaseRoute)
	api.DELETE("/delete/account", security.RateLimiter(security.RequestAmountTypes{
		AuthedUsers: 3,
		GuestUsers:  1,
	}, 60*6), api_delete.DeleteAccountRoute)

	api.GET("/auth/logout", api_auth.AuthLogoutRoute)
	api.POST("/auth/verify", security.RateLimiter(security.RequestAmountTypes{
		AuthedUsers: 3,
		GuestUsers:  3,
	}, 60*1), api_auth.AuthVerifyRoute)
	api.POST("/auth/login", security.RateLimiter(security.RequestAmountTypes{
		AuthedUsers: 1,
		GuestUsers:  10,
	}, 60*1), api_auth.AuthLoginRoute)
	api.POST("/auth/register", security.RateLimiter(security.RequestAmountTypes{
		AuthedUsers: 1,
		GuestUsers:  10,
	}, 60*1), api_auth.AuthRegisterRoute)

	fmt.Printf("zep.run API is running on PORT :: %s\n", os.Getenv("PORT"))
	err := r.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))
	if err != nil {
		return err
	}

	return nil
}
