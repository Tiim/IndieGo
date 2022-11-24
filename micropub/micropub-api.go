package micropub

import (
	"context"
	"fmt"
	"mime/multipart"
	"strings"
	"tiim/go-comment-api/indieauth"

	"github.com/gin-gonic/gin"
)

type micropubApiModule struct {
	store       MicropubStore
	mediaStore  MediaStore
	verifyToken indieauth.TokenVerifier
}

func NewMicropubApiModule(store MicropubStore, mediaStore MediaStore, verifyToken indieauth.TokenVerifier) *micropubApiModule {
	return &micropubApiModule{store: store, mediaStore: mediaStore, verifyToken: verifyToken}
}

func (m *micropubApiModule) Name() string {
	return "micropub"
}

func (m *micropubApiModule) Init(r *gin.Engine) error {
	return nil
}

func (m *micropubApiModule) RegisterRoutes(r *gin.Engine) error {
	r.POST("/micropub", m.micropubEndpoint)
	r.GET("/micropub", m.queryEndpoint)
	return nil
}

func (m *micropubApiModule) micropubEndpoint(c *gin.Context) {
	authorization := c.Request.Header.Get("Authorization")
	authorization = strings.TrimPrefix(authorization, "Bearer ")

	ct := c.ContentType()

	var err error
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
	if data.AccessToken != "" && authorization != "" {
		c.AbortWithError(400, fmt.Errorf("access_token must not be set in both Authorization header and form data"))
		return
	} else if data.AccessToken != "" {
		authorization = data.AccessToken
	}
	scopeChecker, err := m.verifyToken(authorization, []string{"create"})
	if err != nil {
		c.AbortWithError(401, err)
		return
	}

	if data.Action == "create" {
		post := ParseMicropubPost(data)
		if len(data.Files) > 0 {
			err := m.mediaStore.SaveMediaFiles(context.Background(), data, &post)
			if err != nil {
				c.AbortWithError(500, err)
				return
			}
		}
		location, err := m.store.Create(post)
		c.Header("Location", location)
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
		c.Status(202)
	} else if data.Action == "update" {
		if !scopeChecker("update") {
			c.JSON(401, gin.H{"error": "unauthorized"})
			return
		}
		err := m.store.Modify(data.Url, data.Delete, data.Add, data.Replace)
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
		c.Status(200)
	} else if data.Action == "delete" {
		if !scopeChecker("delete") {
			c.JSON(401, gin.H{"error": "unauthorized"})
			return
		}
		err := m.store.Delete(data.Url)
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
		c.Status(200)
	} else {
		c.AbortWithError(400, fmt.Errorf("unsupported action: %s", data.Action))
		return
	}
}

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
	q := c.Query("q")
	if q == "config" {
		c.JSON(200, gin.H{"syndicate-to": []gin.H{}})
	} else if q == "source" {
		url := c.Query("url")
		post, err := m.store.Get(url)
		if err != nil {
			c.AbortWithError(404, err)
			return
		}
		c.JSON(200, post)
	} else if q == "syndicate-to" {
		c.JSON(200, gin.H{"syndicate-to": []gin.H{}})
	} else {
		c.AbortWithError(400, fmt.Errorf("unsupported query: %s", q))
		return
	}
}

func extractFormData(c *gin.Context) (MicropubPostRaw, error) {
	err := c.Request.ParseForm()
	if err != nil {
		return MicropubPostRaw{}, err
	}
	data := make(map[string][]interface{})
	for k, v := range c.Request.Form {
		k = strings.TrimSuffix(k, "[]")
		data[k] = make([]interface{}, len(v))
		for i, s := range v {
			data[k][i] = s
		}
	}
	postType := c.Request.Form["h"]
	if len(postType) == 0 {
		postType = []string{"h-entry"}
	} else {
		postType[0] = "h-" + postType[0]
		delete(data, "h")
	}
	var action string
	actionForm := c.Request.Form["action"]
	if len(actionForm) == 0 {
		action = "create"
	} else {
		delete(data, "action")
		action = actionForm[0]
	}
	var url string
	urlForm := c.Request.Form["url"]
	if len(urlForm) > 0 {
		url = urlForm[0]
		delete(data, "url")
	}
	accessToken := c.Request.FormValue("access_token")
	c.Request.Form.Del("access_token")
	return MicropubPostRaw{
		Action:      action,
		PostTye:     postType,
		Properties:  data,
		AccessToken: accessToken,
		Url:         url,
	}, nil
}

func extractMultipartFormData(c *gin.Context) (MicropubPostRaw, error) {
	err := c.Request.ParseMultipartForm(32 << 20)
	if err != nil {
		return MicropubPostRaw{}, err
	}
	data := make(map[string][]interface{})
	for k, v := range c.Request.MultipartForm.Value {
		k = strings.TrimSuffix(k, "[]")
		data[k] = make([]interface{}, len(v))
		for i, s := range v {
			data[k][i] = s
		}
	}

	files, err := getFiles(c.Request.MultipartForm.File)
	if err != nil {
		return MicropubPostRaw{}, err
	}

	postType := c.Request.MultipartForm.Value["h"]
	if len(postType) == 0 {
		postType = []string{"h-entry"}
	} else {
		postType[0] = "h-" + postType[0]
		delete(data, "h")
	}
	var action string
	actionForm := c.Request.MultipartForm.Value["action"]
	if len(actionForm) == 0 {
		action = "create"
	} else {
		delete(data, "action")
		action = actionForm[0]
	}
	var accessToken string
	accessTokenForm := c.Request.MultipartForm.Value["access_token"]
	if len(accessToken) > 0 {
		accessToken = accessTokenForm[0]
		delete(data, "access_token")
	}
	var url string
	urlForm := c.Request.MultipartForm.Value["url"]
	if len(urlForm) > 0 {
		url = urlForm[0]
		delete(data, "url")
	}
	return MicropubPostRaw{
		Action:      action,
		PostTye:     postType,
		Properties:  data,
		AccessToken: accessToken,
		Url:         url,
		Files:       files,
	}, nil
}

func extractJSONData(c *gin.Context) (MicropubPostRaw, error) {
	var data MicropubPostRaw
	err := c.BindJSON(&data)

	if err != nil {
		return MicropubPostRaw{}, err
	}
	if data.Action == "" {
		data.Action = "create"
	}
	return data, nil
}

func getFiles(files map[string][]*multipart.FileHeader) ([]MicropubFile, error) {
	mpfiles := make([]MicropubFile, 0)
	for _, v := range files {
		for _, f := range v {
			file, err := f.Open()
			if err != nil {
				return nil, err
			}
			mpfiles = append(mpfiles, MicropubFile{
				Name:        f.Filename,
				ContentType: f.Header.Get("Content-Type"),
				Size:        f.Size,
				Reader:      file,
			})
		}
	}
	return mpfiles, nil
}
