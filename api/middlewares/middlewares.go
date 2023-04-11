package middlewares

import (
	"net/http"
	"strings"

	"CRUD-API/handlers/auth"

	"github.com/gin-gonic/gin"
)

func Authorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		h, ok := c.Request.Header["Authorization"]

		if !ok {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		term := strings.Fields(h[0])
		method, token := term[0], term[1]

		if !strings.EqualFold("Bearer", method) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if _, err := auth.ValidateToken(token); err != nil {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"message": err.Error(),
				},
			)
			return
		}

		c.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// c.Writer.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5173")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
