package micropub

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

func (m *micropubApiModule) queryEndpoint(c *gin.Context) {
	authorization := c.Request.Header.Get("Authorization")
	authorization = strings.TrimPrefix(authorization, "Bearer ")
	authQuery := c.Query("access_token")
	if authQuery != "" && authorization != "" {
		c.AbortWithError(400, fmt.Errorf("access_token must not be set in both Authorization header and form data"))
		return
	} else if authQuery != "" {
		authorization = authQuery
	}
	_, err := m.verifyToken(authorization, []string{"create"})
	if err != nil {
		c.AbortWithError(401, err)
		return
	}

	switch c.Query("q") {
	case "config", "syndicate-to":
		c.JSON(200, gin.H{
			"media-endpoint": "/micropub/media",
			"syndicate-to":   []string{},
		})
	case "source":
		url := c.Query("url")
		post, err := m.store.Get(url)
		if err != nil {
			c.AbortWithError(404, err)
			return
		}
		c.JSON(200, post)
	default:
		c.AbortWithError(400, fmt.Errorf("unsupported query: %s", c.Query("q")))
		return
	}
}
