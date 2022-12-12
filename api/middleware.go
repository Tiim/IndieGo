package api

// request_logger.go
import (
	"strings"

	"github.com/gin-gonic/gin"
)

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, POST, PUT, DELETE, OPTIONS, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func trailingSlash(e *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if path != "/" && strings.HasSuffix(path, "/") {
			c.Request.URL.Path = strings.TrimSuffix(path, "/")
			e.HandleContext(c)
			c.Abort()
		} else {
			c.Next()
		}
	}
}
