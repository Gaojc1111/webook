package middlewares

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// LoginMiddlewareBuilder 整合登录相关middleware
type LoginMiddlewareBuilder struct {
	paths []string
}

// IgnorePaths 对不用身份校验的HTTP请求放行
// 返回*LoginMiddlewareBuilder的意义是：可以连续调用IgnorePaths
func (l *LoginMiddlewareBuilder) IgnorePaths(path string) *LoginMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}

func (l *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 整合放行路径判断
		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
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
