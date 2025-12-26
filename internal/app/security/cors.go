package security

import (
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORSHandler() gin.HandlerFunc {
	apiUrl := os.Getenv("API_URL")
	return cors.New(cors.Config{
		AllowOrigins:     []string{apiUrl},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "DELETE"},
		AllowCredentials: true,
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowOriginFunc: func(origin string) bool {
			return origin == apiUrl
		},
		MaxAge: 12 * time.Hour,
	})
}
