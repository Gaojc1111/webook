package ioc

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
	"webook/internal/web"
	"webook/internal/web/middlewares"
	"webook/pkg/ginx/middleware/ratelimit"
	"webook/pkg/limiter"
)

func InitWebServer(userHandler *web.UserHandler, wechatHandler *web.OAuth2WechatHandler, middlewares []gin.HandlerFunc) *gin.Engine {
	server := gin.Default()
	server.Use(middlewares...)
	wechatHandler.RegisterRoutes(server)
	userHandler.RegisterRoutes(server)
	return server
}

func InitGinMiddlewares(redisClient redis.Cmdable) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		// 中间件 先注册先执行
		// 解决跨域问题
		// https://github.com/gin-contrib/cors
		cors.New(cors.Config{
			//AllowOrigins:     []string{"http://localhost:3000"},
			AllowMethods: []string{"PUT", "PATCH", "GET", "POST"},
			AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
			// JWT 放行
			ExposeHeaders: []string{"Content-Length", "x-jwt-token", "x-refresh-token"},

			AllowCredentials: true,
			// 放行所有包含http://localhost 前缀的域名
			AllowOriginFunc: func(origin string) bool {
				return strings.HasPrefix(origin, "http://localhost")
			},
			MaxAge: 12 * time.Hour,
		}),
		// jwt 中间件
		middlewares.NewLoginJWTMiddlewareBuilder().
			IgnorePaths("/users/login").
			IgnorePaths("/users/signup").
			IgnorePaths("/users/login_sms/code/send").
			IgnorePaths("/users/login_sms").
			IgnorePaths("/oauth2/wechat/authurl").
			IgnorePaths("/oauth2/wechat/callback").
			Build(),
		// redis限流中间件
		ratelimit.NewBuilder(limiter.NewRedisSlidingWindowLimiter(redisClient, time.Second, 1000)).Build(),
		gin.Logger(),
		gin.Recovery(),
	}
}
