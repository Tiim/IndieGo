package api

// request_logger.go
import (
	"log"

	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, POST, PUT, DELETE, OPTIONS, PATCH")

		log.Println("CORS", c.Request.Method, c.Request.URL.Path)

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
		if path != "/" && path[len(path)-1] == '/' {
			c.Request.URL.Path = path[:len(path)-1]
			e.HandleContext(c)
			c.Abort()
		} else {
			c.Next()
		}
	}
}
