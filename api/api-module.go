package api

import "github.com/gin-gonic/gin"

type ApiModule interface {
	Name() string
	Init(r *gin.Engine) error
	RegisterRoutes(r *gin.Engine) error
}
