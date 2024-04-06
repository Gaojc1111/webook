package middlewares

import (
	"LittleRedBook/internal/web"
	"encoding/gob"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"time"
)

// LoginJWTMiddlewareBuilder JWT登录校验
type LoginJWTMiddlewareBuilder struct {
	paths []string
}

func NewLoginJWTMiddlewareBuilder() *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{}
}

// IgnorePaths 对不用身份校验的HTTP请求放行
// 返回*LoginMiddlewareBuilder的意义是：可以连续调用IgnorePaths
func (l *LoginJWTMiddlewareBuilder) IgnorePaths(path string) *LoginJWTMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
	// 用Go的方式编码解码
	gob.Register(time.Now())
	return func(ctx *gin.Context) {
		// 整合放行路径判断
		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}

		// JWT校验
		tokenStr := ctx.GetHeader("Authorization")
		if tokenStr == "" {
			// 没登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		segs := strings.Split(tokenStr, " ") // Bearer token...
		if len(segs) != 2 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr = segs[1]
		claims := &web.UserClaims{}
		// token校验，ParseWithClaims要传指针
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("Hbzhtd0211"), nil
		})
		if err != nil {
			// 没登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// err为nil, token不为nil 约定
		if !token.Valid {
			// 没登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if claims.UserAgent != ctx.Request.UserAgent() {
			// 比如： 登录在谷歌，其他操作在bing
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if claims.UserID == 0 {
			// 没登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}

	}
}
