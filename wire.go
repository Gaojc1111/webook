//go:build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"webook/internal/repository"
	"webook/internal/repository/cache"
	"webook/internal/repository/dao"
	"webook/internal/service"
	"webook/internal/web"
	"webook/ioc"
)

func initWebServer() *gin.Engine {
	wire.Build(
		// 底层存储
		ioc.InitDB, ioc.InitRedis,

		// dao & cache
		dao.NewUserDAO,
		cache.NewUserCache, cache.NewCodeCache,

		// repository
		repository.NewCachedUserRepository, repository.NewCodeRepository,

		// service
		ioc.InitSMSService, ioc.InitWechatService,
		service.NewUserService, service.NewCodeService,

		// handler
		web.NewUserHandler, ioc.InitGinMiddlewares, ioc.InitWebServer,
		web.NewOAuth2WechatHandler,
	)
	return gin.Default()
}
