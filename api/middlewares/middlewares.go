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

var whiteList = map[string]bool{
	"http://localhost:3000": true,
}

// func CORSMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		if origin := c.Request.Header.Get("Origin"); whiteList[origin] {
// 			// c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
// 			c.Header("Access-Control-Allow-Origin", origin)
// 		}
// 		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
// 		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")
// 		if c.Request.Method == "OPTIONS" {
// 			c.AbortWithStatus(204)
// 			return
// 		}
// 		c.Next()
// 	}
// }

// Jacky
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE,PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
