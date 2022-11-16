package indieauth

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strings"

	_ "embed"

	"github.com/gin-gonic/gin"
)

//go:embed authorize.tmpl
var authorizeTemplate string

type indieAuthApiModule struct {
	store               Store
	httpClient          http.Client
	group               *gin.RouterGroup
	baseUrl             string
	profileCanonicalUrl string
	password            string
	authorizeTemplate   *template.Template
}

func NewIndieAuthApiModule(baseUrl, profileCanonicalUrl, password string, store Store, client http.Client) *indieAuthApiModule {
	authorizeTemplate := template.Must(template.New("authorize").Parse(authorizeTemplate))
	return &indieAuthApiModule{
		profileCanonicalUrl: profileCanonicalUrl,
		baseUrl:             baseUrl,
		authorizeTemplate:   authorizeTemplate,
		httpClient:          client,
		store:               store,
		password:            password,
	}
}

func (m *indieAuthApiModule) Name() string {
	return "IndieAuth"
}

func (m *indieAuthApiModule) Init(r *gin.Engine) error {
	m.group = r.Group("/indieauth")
	return nil
}

func (m *indieAuthApiModule) RegisterRoutes(r *gin.Engine) error {
	m.group.GET("/metadata", m.metadataEndpoint)
	m.group.POST("/token", m.tokenEndpoint)
	m.group.GET("/authorize", m.authorizeEndpoint)
	m.group.POST("/authorize", m.tokenEndpoint)
	m.group.POST("/introspection", m.introspectionEndpoint)
	m.group.POST("/login", m.loginEndpoint)
	return nil
}

func (m *indieAuthApiModule) metadataEndpoint(c *gin.Context) {

	challengeNames := make([]string, 0, len(challenges))
	for name := range challenges {
		challengeNames = append(challengeNames, name)
	}

	c.JSON(200, gin.H{
		"issuer":                           m.baseUrl + "/indieauth",
		"authorization_endpoint":           m.baseUrl + "/indieauth/authorize",
		"token_endpoint":                   m.baseUrl + "/indieauth/token",
		"introspection_endpoint ":          m.baseUrl + "/indieauth/introspection",
		"scopes_supported":                 scopes,
		"code_challenge_methods_supported": challengeNames,
	})
}

func (m *indieAuthApiModule) authorizeEndpoint(c *gin.Context) {
	responseType := c.Query("response_type")
	if responseType != "code" {
		c.AbortWithError(400, fmt.Errorf("invalid response_type, only 'code' is supported"))
		return
	}
	clientId := c.Query("client_id")
	redirectUri := c.Query("redirect_uri")
	state := c.Query("state")
	codeChallenge := c.Query("code_challenge")
	codeChallengeMethod := c.Query("code_challenge_method")
	scope := c.Query("scope")
	me := c.Query("me")

	warnings := make([]string, 0)

	_, ok := challenges[codeChallengeMethod]
	if !ok {
		c.AbortWithError(400, fmt.Errorf("invalid code_challenge_method"))
		return
	}
	if scope == "" {
		warnings = append(warnings, "no scope specified, using default scope")
		scope = "profile"
	}

	code, err := newAuthCode(redirectUri, clientId, scope, state, codeChallenge, codeChallengeMethod, me)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	appInfo, err := getAppInfo(clientId, m.httpClient, c.Request.Context())
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	if !strings.HasPrefix(redirectUri, clientId) {
		redirectOk := false
		for _, redirect := range appInfo.RedirectUris {
			if redirectUri == redirect {
				redirectOk = true
				break
			}
		}
		if !redirectOk {
			c.AbortWithError(400, fmt.Errorf("invalid redirect_uri"))
			return
		}
	}

	err = m.store.StoreAuthCode(code)

	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.Header("Content-Type", "text/html")
	m.authorizeTemplate.Execute(c.Writer, gin.H{"Code": code.code, "AppInfo": appInfo, "Warnings": warnings, "Scope": scope})
}

func (m *indieAuthApiModule) tokenEndpoint(c *gin.Context) {
	grantType := c.Request.FormValue("grant_type")
	code := c.Request.FormValue("code")
	clientId := c.Request.FormValue("client_id")
	redirectUri := c.Request.FormValue("redirect_uri")
	codeVerifier := c.Request.FormValue("code_verifier")

	if grantType != "authorization_code" {
		c.AbortWithError(400, fmt.Errorf("invalid grant_type, only 'authorization_code' is supported"))
		return
	}

	authCode, err := m.store.GetAuthCode(code)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	if authCode == nil {
		c.AbortWithError(400, fmt.Errorf("invalid code"))
		return
	}

	if authCode.clientId != clientId {
		c.AbortWithError(400, fmt.Errorf("invalid client_id"))
		return
	}
	if authCode.redirectUri != redirectUri {
		c.AbortWithError(400, fmt.Errorf("invalid redirect_uri"))
		return
	}

	challenge := challenges[authCode.codeChallengeMethod]
	if !challenge.Verify(codeVerifier, authCode.codeChallenge) {
		c.AbortWithError(400, fmt.Errorf("invalid code_verifier"))
		return
	}
	if c.Request.URL.Path == "/indieauth/authorize" {
		// only a profile request, do not issue a token
		c.JSON(200, gin.H{"me": m.profileCanonicalUrl})
	} else {
		c.JSON(200, gin.H{"me": m.profileCanonicalUrl})
	}
}

func (m *indieAuthApiModule) introspectionEndpoint(c *gin.Context) {
}

func (m *indieAuthApiModule) loginEndpoint(c *gin.Context) {
	code := c.Request.FormValue("code")
	password := c.Request.FormValue("password")

	if password != m.password {
		c.AbortWithError(401, fmt.Errorf("invalid password"))
		return
	}

	authCode, err := m.store.GetAuthCode(code)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	queryValues := url.Values{}
	queryValues.Set("code", authCode.code)
	queryValues.Set("state", authCode.state)
	queryValues.Set("iss", m.baseUrl+"/indieauth")

	c.Redirect(302, authCode.redirectUri+"?"+queryValues.Encode())
}