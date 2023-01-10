package micropub

import (
	"fmt"
	"log"
	"tiim/go-comment-api/plugins/indieauth"

	"github.com/gin-gonic/gin"
)

func (m *micropubApiModule) micropubEndpoint(c *gin.Context) {
	authorization, err := authToken(c)
	if err != nil {
		c.AbortWithError(401, err)
		return
	}
	scopeChecker, err := m.verifyToken(authorization, []string{"create"})
	if err != nil {
		c.AbortWithError(401, err)
		return
	}

	ct := c.ContentType()
	var data MicropubPostRaw
	if ct == "application/x-www-form-urlencoded" {
		data, err = extractFormData(c)
	} else if ct == "multipart/form-data" {
		data, err = extractMultipartFormData(c)
	} else if ct == "application/json" {
		data, err = extractJSONData(c)
	} else {
		c.AbortWithError(400, fmt.Errorf("unsupported Content-Type: %s", ct))
		return
	}
	if err != nil {
		c.AbortWithError(400, fmt.Errorf("failed to parse request: %w", err))
		return
	}

	switch data.Action {
	case "create":
		m.actionCreate(c, data)
	case "update":
		m.actionUpdate(c, data, scopeChecker)
	case "delete":
		m.actionDelete(c, data, scopeChecker)
	default:
		c.AbortWithError(400, fmt.Errorf("unsupported action: %s", data.Action))
		return
	}
}

func (m *micropubApiModule) actionCreate(c *gin.Context, data MicropubPostRaw) {
	post := ParseMicropubPost(data)

	for _, file := range data.Files {
		url, err := m.mediaStore.SaveMediaFiles(c, file)
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
		addUrlToPost(&post, url, file.Name, file.ContentType, m.logger)
	}
	location, err := m.store.Create(post)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	log.Printf("Created post at %s", location)
	c.Header("Location", location)
	c.Status(202)
}

func (m *micropubApiModule) actionUpdate(c *gin.Context, data MicropubPostRaw, scopeChecker indieauth.ScopeCheck) {
	if !scopeChecker("update") {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	err := m.store.Modify(data.Url, data.Delete, data.Add, data.Replace)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	log.Printf("Updated post %s", data.Url)
	c.Status(200)
}

func (m *micropubApiModule) actionDelete(c *gin.Context, data MicropubPostRaw, scopeChecker indieauth.ScopeCheck) {
	if !scopeChecker("delete") {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	err := m.store.Delete(data.Url)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	log.Printf("Deleted post %s", data.Url)
	c.Status(200)
}
