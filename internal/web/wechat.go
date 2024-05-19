package web

import (
	"github.com/gin-gonic/gin"
	"webook/internal/service"
	"webook/internal/service/oauth2/wechat"
)

type OAuth2WechatHandler struct {
	svc     wechat.Service
	userSvc service.UserService
	JWTHandler
}

func NewOAuth2WechatHandler(svc wechat.Service) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		svc: svc,
	}
}

func (o *OAuth2WechatHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/oauth2/wechat")
	g.GET("/authurl", o.Auth2URL)
	g.Any("/callback", o.CallBack)
}

func (o *OAuth2WechatHandler) Auth2URL(ctx *gin.Context) {
	url, err := o.svc.AuthURL(ctx)
	if err != nil {
		ctx.JSON(200, Result{
			Code: 5,
			Msg:  "URL跳转失败",
		})
		return
	}
	ctx.JSON(200, Result{
		Data: url,
	})
}

func (o *OAuth2WechatHandler) CallBack(ctx *gin.Context) {
	code := ctx.Query("code")
	//state := ctx.Query("state")
	wechatInfo, err := o.svc.VerifyCode(ctx, code)
	if err != nil {
		ctx.JSON(200, Result{
			Code: 4,
			Msg:  "微信授权失败",
		})
		return
	}
	ctx.JSON(200, Result{
		Data: wechatInfo,
	})
	user, err := o.userSvc.FindOrCreateByWechat(ctx, wechatInfo)
	if err != nil {
		ctx.JSON(200, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	o.setJWTToken(ctx, user.ID)
	ctx.JSON(200, Result{
		Msg: "ok",
	})
}
