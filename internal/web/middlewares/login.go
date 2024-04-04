package middlewares

import (
	"encoding/gob"
	"net/http"
	"time"

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
	// 用Go的方式编码解码
	gob.Register(time.Now())
	return func(ctx *gin.Context) {
		// 整合放行路径判断
		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}

		// session校验
		session := sessions.Default(ctx)
		userID := session.Get("userID")
		if userID == nil {
			//未登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		//刷新session
		session.Set("userID", userID)
		session.Options(sessions.Options{
			MaxAge: 60 * 30, //30min
		})
		//
		now := time.Now()
		updateTime := session.Get("update_time")
		// 还没刷新过
		if updateTime == nil {
			session.Set("updateTime", now)
			err := session.Save()
			if err != nil {
				panic(err)
			}
			return
		}

		updateTimeVal, ok := updateTime.(time.Time)
		if !ok {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// 已经刷新过，判断是否要续期
		if now.Sub(updateTimeVal) > time.Minute {
			session.Set("update_time", now)
			err := session.Save()
			if err != nil {
				panic(err)
			}
		}
	}
}
