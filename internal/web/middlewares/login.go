package middlewares

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// LoginMiddlewareBuilder 整合登录相关middleware
type LoginMiddlewareBuilder struct {
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}

func (l *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.URL.Path == "/users/login" {
			return
		}
		if ctx.Request.URL.Path == "/users/signup" {
			return
		}
		session := sessions.Default(ctx)
		userID := session.Get("userID")
		if userID == nil {
			//未登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

	}
}
