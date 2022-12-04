package indieauth

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strings"
	"tiim/go-comment-api/config"

	_ "embed"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v4"
)

//go:embed authorize.tmpl
var authorizeTemplate string

// TODO: find proper names for the plugin instances
type IndieAuthApiModule struct {
	store               Store
	httpClient          http.Client
	group               *gin.RouterGroup
	baseUrl             string
	profileCanonicalUrl string
	password            string
	jwtSecret           string
	authorizeTemplate   *template.Template
}

func NewIndieAuthApiModule(baseUrl, profileCanonicalUrl, password, jwtSecret string, store Store, client http.Client) *IndieAuthApiModule {
	authorizeTemplate := template.Must(template.New("authorize").Parse(authorizeTemplate))
	return &IndieAuthApiModule{
		profileCanonicalUrl: profileCanonicalUrl,
		baseUrl:             baseUrl,
		authorizeTemplate:   authorizeTemplate,
		httpClient:          client,
		store:               store,
		password:            password,
		jwtSecret:           jwtSecret,
	}
}

func (m *IndieAuthApiModule) Name() string {
	return "indieauth"
}

func (m *IndieAuthApiModule) Start() error {
	return nil
}

func (m *IndieAuthApiModule) Init(config config.GlobalConfig) error {
	return nil
}

func (m *IndieAuthApiModule) InitGroups(r *gin.Engine) error {
	m.group = r.Group("/indieauth")
	return nil
}

func (m *IndieAuthApiModule) RegisterRoutes(r *gin.Engine) error {
	m.group.GET("/metadata", m.metadataEndpoint)
	m.group.POST("/token", m.tokenEndpoint)
	m.group.GET("/token", m.introspectionEndpoint)
	m.group.GET("/authorize", m.authorizeEndpoint)
	m.group.POST("/authorize", m.tokenEndpoint)
	m.group.POST("/introspection", m.introspectionEndpoint)
	m.group.POST("/login", m.loginEndpoint)
	return nil
}

func (m *IndieAuthApiModule) metadataEndpoint(c *gin.Context) {

	challengeNames := make([]string, 0, len(challenges))
	for name := range challenges {
		challengeNames = append(challengeNames, name)
	}

	c.JSON(200, gin.H{
		"issuer":                           m.baseUrl + "/indieauth",
		"authorization_endpoint":           m.baseUrl + "/indieauth/authorize",
		"token_endpoint":                   m.baseUrl + "/indieauth/token",
		"introspection_endpoint ":          m.baseUrl + "/indieauth/introspection",
		"code_challenge_methods_supported": challengeNames,
	})
}

func (m *IndieAuthApiModule) authorizeEndpoint(c *gin.Context) {
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

	if codeChallengeMethod == "" {
		codeChallengeMethod = "plain"
	}
	if codeChallengeMethod == "plain" {
		warnings = append(warnings, "using insecure code_challenge_method 'plain'")
	}
	_, ok := challenges[codeChallengeMethod]
	if !ok {
		c.AbortWithError(400, fmt.Errorf("invalid code_challenge_method %s", codeChallengeMethod))
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
	m.authorizeTemplate.Execute(c.Writer, gin.H{"Code": code.code, "AppInfo": appInfo, "Warnings": warnings, "Scopes": strings.Split(scope, " "), "Me": me})
}

func (m *IndieAuthApiModule) tokenEndpoint(c *gin.Context) {
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
		c.AbortWithError(400, fmt.Errorf("invalid code_verifier: %s code_challenge %s", codeVerifier, authCode.codeChallenge))
		return
	}

	accessToken, err := m.store.RedeemAccessToken(code)

	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	if c.Request.URL.Path == "/indieauth/authorize" {
		// only a profile request, do not issue a token
		c.JSON(200, gin.H{"me": m.profileCanonicalUrl})
	} else {

		jwtClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"iss": m.baseUrl,
			"sub": m.profileCanonicalUrl,
			"aud": accessToken.clientId,
			"iat": accessToken.issuedAt.Unix(),
			"exp": accessToken.expiresAt.Unix(),
			m.baseUrl: map[string]interface{}{
				"scope": accessToken.scope,
			},
		})

		token, error := jwtClaim.SignedString([]byte(m.jwtSecret))

		if error != nil {
			c.AbortWithError(500, error)
			return
		}

		c.JSON(200, gin.H{
			"me":           m.profileCanonicalUrl,
			"scope":        accessToken.scope,
			"access_token": token,
			"token_type":   "bearer",
		})
	}
}

func (m *IndieAuthApiModule) introspectionEndpoint(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		authHeader = strings.TrimPrefix(authHeader, "Bearer ")
	} else {
		authHeader = ""
	}
	tokenString := c.Request.FormValue("token")
	if tokenString == "" {
		tokenString = authHeader
	}
	if tokenString == "" {
		c.AbortWithError(400, fmt.Errorf("no token provided"))
		return
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.jwtSecret), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid && err == nil {
		c.JSON(200, gin.H{
			"active":    true,
			"me":        claims["sub"],
			"scope":     claims[m.baseUrl].(map[string]interface{})["scope"],
			"client_id": claims["aud"],
			"exp":       claims["exp"],
			"iat":       claims["iat"],
		})
	} else {
		c.JSON(200, gin.H{"active": false})
	}
}

func (m *IndieAuthApiModule) VerifyToken(tokenString string, minimalScopes []string) (ScopeCheck, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && err == nil && token.Valid {
		// check scopes
		scopes := strings.Split(claims[m.baseUrl].(map[string]interface{})["scope"].(string), " ")
		for _, minimalScope := range minimalScopes {
			found := false
			for _, scope := range scopes {
				if scope == minimalScope {
					found = true
					break
				}
			}
			if !found {
				return nil, fmt.Errorf("missing scope %s", minimalScope)
			}
		}

		scopeChecker := func(scope string) bool {
			for _, s := range scopes {
				if s == scope {
					return true
				}
			}
			return false
		}

		return scopeChecker, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}

func (m *IndieAuthApiModule) loginEndpoint(c *gin.Context) {
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

	approvedScopes := make([]string, 0)
	for _, scope := range strings.Split(authCode.scope, " ") {
		if c.Request.FormValue("scope-"+scope) == "true" {
			approvedScopes = append(approvedScopes, scope)
		}
	}
	m.store.UpdateScope(code, strings.Join(approvedScopes, " "))

	queryValues := url.Values{}
	queryValues.Set("code", authCode.code)
	queryValues.Set("state", authCode.state)
	queryValues.Set("iss", m.baseUrl+"/indieauth")

	c.Redirect(302, authCode.redirectUri+"?"+queryValues.Encode())
}
