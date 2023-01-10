package micropub

import (
	"fmt"
	"log"
	"mime/multipart"
	"strings"
	"tiim/go-comment-api/config"
	"tiim/go-comment-api/lib/mfobjects"
	"tiim/go-comment-api/plugins/indieauth"

	"github.com/gin-gonic/gin"
)

type micropubApiModule struct {
	store       micropubStore
	mediaStore  mediaStore
	verifyToken indieauth.TokenVerifier
	logger      *log.Logger
}

func newMicropubApiModule(store micropubStore, mediaStore mediaStore, verifyToken indieauth.TokenVerifier, logger *log.Logger) *micropubApiModule {
	return &micropubApiModule{store: store, mediaStore: mediaStore, verifyToken: verifyToken, logger: logger}
}

func (m *micropubApiModule) Name() string {
	return "micropub"
}

func (m *micropubApiModule) Init(config config.GlobalConfig) error {
	return nil
}

func (m *micropubApiModule) RegisterRoutes(r *gin.Engine) error {
	r.POST("/micropub", m.micropubEndpoint)
	r.GET("/micropub", m.queryEndpoint)
	r.POST("/micropub/media", m.mediaEndpoint)
	return nil
}

func (m *micropubApiModule) Start() error {
	return nil
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
	c.Request.Form.Del("access_token")
	return MicropubPostRaw{
		Action:     action,
		PostTye:    postType,
		Properties: data,
		Url:        url,
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
	var url string
	urlForm := c.Request.MultipartForm.Value["url"]
	if len(urlForm) > 0 {
		url = urlForm[0]
		delete(data, "url")
	}
	return MicropubPostRaw{
		Action:     action,
		PostTye:    postType,
		Properties: data,
		Url:        url,
		Files:      files,
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

func addUrlToPost(mp *MicropubPost, url, name, contentType string, logger *log.Logger) {
	switch contentType {
	case "image/jpeg", "image/png", "image/gif":
		mp.Entry.Photos = append(mp.Entry.Photos, mfobjects.MF2Photo{
			Url: url,
		})
	default:
		logger.Println("Unknown content type to add to MicropubPost: ", contentType)
	}
}

func authToken(c *gin.Context) (string, error) {
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader != "" {
		authHeader = strings.TrimPrefix(authHeader, "Bearer ")
	}
	authForm := c.Request.FormValue("access_token")
	if authForm != "" {
		c.Request.Form.Del("access_token")
	}

	if authHeader != "" && authForm != "" {
		return "", fmt.Errorf("both authorization header and access_token form value are set")
	}
	if authHeader != "" {
		return authHeader, nil
	} else if authForm != "" {
		return authForm, nil
	}
	return "", fmt.Errorf("no authorization header or access_token form value set")
}
