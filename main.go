package main

import (
	"strings"
	"time"
	"webook/config"
	"webook/internal/repository"
	"webook/internal/repository/dao"
	"webook/internal/service"
	"webook/internal/web"
	"webook/internal/web/middlewares"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db := initDB()
	server := initWebserver()

	u := initUser(db)
	u.RegisterRoutes(server)
	//pprof.Register(server) // 性能分析: 注册pprof相关路由
	//server := gin.Default()
	err := server.Run(":8080")
	if err != nil {
		return
	}
}

func initWebserver() *gin.Engine {
	server := gin.Default()

	// 中间件 先注册先执行
	// 解决跨域问题
	// https://github.com/gin-contrib/cors
	server.Use(cors.New(cors.Config{
		//AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods: []string{"PUT", "PATCH", "GET", "POST"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
		// JWT 放行
		ExposeHeaders: []string{"Content-Length", "x-jwt-token"},

		AllowCredentials: true,
		// 放行所有包含http://localhost 前缀的域名
		AllowOriginFunc: func(origin string) bool {
			return strings.HasPrefix(origin, "http://localhost")
		},
		MaxAge: 12 * time.Hour,
	}))

	// redis限流
	//redisClient := redis.NewClient(&redis.Options{
	//	Addr: "localhost:6379",
	//})
	//
	//server.Use(ratelimit.NewBuilder(redisClient,
	//	time.Second, 1).Build())

	// JWT 验证
	server.Use(middlewares.NewLoginJWTMiddlewareBuilder().
		IgnorePaths("/users/login").
		IgnorePaths("/users/signup").
		Build())

	return server
}

func initUser(db *gorm.DB) *web.UserHandler {
	ud := dao.NewUserDAO(db)
	repo := repository.NewUserRepository(ud)
	svc := service.NewUserService(repo)
	u := web.NewUserHandler(svc)
	return u
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		panic(err)
	}
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}
