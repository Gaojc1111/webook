package main

import (
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
	"webook/config"
	"webook/internal/repository"
	"webook/internal/repository/cache"
	"webook/internal/repository/dao"
	"webook/internal/service"
	"webook/internal/service/sms/localsms"
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
		IgnorePaths("/users/login_sms/code/send").
		IgnorePaths("/users/login_sms").
		Build())

	return server
}

func initUser(db *gorm.DB) *web.UserHandler {
	ud := dao.NewUserDAO(db)
	client := redis.NewClient(&redis.Options{
		Addr:     config.Config.Redis.Addr,
		Password: "",
		DB:       0,
	})
	//codeSv
	user_cache := cache.NewUserCache(client)
	code_cache := cache.NewCodeCache(client)
	user_repo := repository.NewUserRepository(ud, user_cache)
	code_repo := repository.NewCodeRepository(code_cache)
	userSvc := service.NewUserService(user_repo)
	codeSvc := service.NewCodeService(code_repo, localsms.NewService())
	userHandler := web.NewUserHandler(userSvc, codeSvc)
	return userHandler
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
