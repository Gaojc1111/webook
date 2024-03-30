package main

import (
	"Learn/LittleRedBook/internal/repository"
	"Learn/LittleRedBook/internal/repository/dao"
	"Learn/LittleRedBook/internal/service"
	"Learn/LittleRedBook/internal/web"
	"strings"
	"time"

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

	server.Run(":8080")
}

func initWebserver() *gin.Engine {
	server := gin.Default()

	// 中间件 先注册先执行
	// 解决跨域问题
	// https://github.com/gin-contrib/cors
	server.Use(cors.New(cors.Config{
		//AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:  []string{"PUT", "PATCH", "GET", "POST"},
		AllowHeaders:  []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders: []string{"Content-Length"},

		AllowCredentials: true,
		// 放行所有包含http://localhost 前缀的域名
		AllowOriginFunc: func(origin string) bool {
			return strings.HasPrefix(origin, "http://localhost")
		},
		MaxAge: 12 * time.Hour,
	}))
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
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:3306)/RedBook"))
	if err != nil {
		panic(err)
	}
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}
