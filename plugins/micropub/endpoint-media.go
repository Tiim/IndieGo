package micropub

import (
	"context"

	"github.com/gin-gonic/gin"
)

func (m *micropubApiModule) mediaEndpoint(c *gin.Context) {
	authorization, err := authToken(c)
	if err != nil {
		c.AbortWithError(401, err)
		return
	}
	_, err = m.verifyToken(authorization, []string{"create"})
	if err != nil {
		c.AbortWithError(401, err)
		return
	}

	f, err := c.FormFile("file")
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	file, err := f.Open()
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	mpFile := MicropubFile{
		Name:        f.Filename,
		ContentType: f.Header.Get("Content-Type"),
		Size:        f.Size,
		Reader:      file,
	}
	url, err := m.mediaStore.SaveMediaFiles(context.Background(), mpFile)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	c.Header("Location", url)
}
