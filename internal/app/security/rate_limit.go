package security

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type RequestAmountTypes struct {
	AuthedUsers int
	GuestUsers  int
}

func RateLimiter(requestAmountTypes RequestAmountTypes, secondLength int) gin.HandlerFunc {
	authLimit := rate.Limit(float64(requestAmountTypes.AuthedUsers) / float64(secondLength))
	authLimiter := rate.NewLimiter(authLimit, requestAmountTypes.AuthedUsers)

	guestLimit := rate.Limit(float64(requestAmountTypes.GuestUsers) / float64(secondLength))
	guestLimiter := rate.NewLimiter(guestLimit, requestAmountTypes.GuestUsers)

	return func(c *gin.Context) {
		user, x := c.Get("user")
		path := c.FullPath()
		if path == "/api/v2/post/screenshot" && x {
			userMap := user.(gin.H)
			userIsSubscribed := userMap["is_subscribed"].(bool)
			if userIsSubscribed {
				c.Next()
				return
			}
		}

		var limiter *rate.Limiter
		if x {
			limiter = authLimiter
		} else if !x {
			limiter = guestLimiter
		}

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"error":   "__REQUEST__",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
