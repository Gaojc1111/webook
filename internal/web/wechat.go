package web

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/lithammer/shortuuid/v4"
	"net/http"
	"webook/internal/service"
	"webook/internal/service/oauth2/wechat"
)

type OAuth2WechatHandler struct {
	svc     wechat.Service
	userSvc service.UserService
	JWTHandler
	key             []byte
	stateCookieName string
}

type StateClaims struct {
	jwt.RegisteredClaims
	State string
}

func NewOAuth2WechatHandler(svc wechat.Service) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		svc:             svc,
		key:             []byte("Hbzhtd0211"),
		stateCookieName: "jwt_state",
		JWTHandler:      NewJWTHandler(),
	}

}

func (o *OAuth2WechatHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/oauth2/wechat")
	g.GET("/authurl", o.Auth2URL)
	g.Any("/callback", o.CallBack)
}

func (o *OAuth2WechatHandler) Auth2URL(ctx *gin.Context) {
	state := uuid.New()
	url, err := o.svc.AuthURL(ctx, state)
	if err != nil {
		ctx.JSON(200, Result{
			Code: 5,
			Msg:  "URL跳转失败",
		})
		return
	}
	err = o.setStateCookie(ctx, state)
	if err != nil {
		ctx.JSON(200, Result{
			Code: 4,
			Msg:  "系统错误",
		})
		return
	}
	ctx.JSON(200, Result{
		Data: url,
	})
}

func (o *OAuth2WechatHandler) CallBack(ctx *gin.Context) {
	err := o.VerifyState(ctx)
	if err != nil {
		ctx.JSON(200, Result{
			Code: 3,
			Msg:  "非法请求",
		})
		return
	}
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

func (o *OAuth2WechatHandler) VerifyState(ctx *gin.Context) error {
	state := ctx.Query("state")
	stateCookie, err := ctx.Cookie(o.stateCookieName)
	if err != nil {
		return fmt.Errorf("state cookie 不存在: %s", err)
	}
	var claims StateClaims
	_, err = jwt.ParseWithClaims(stateCookie, &claims, func(token *jwt.Token) (interface{}, error) {
		return o.key, nil
	})
	if err != nil {
		return fmt.Errorf("state cookie 无效: %s", err)
	}
	if state != claims.State {
		return errors.New("state 不匹配")
	}

	return nil
}

func (o *OAuth2WechatHandler) setStateCookie(ctx *gin.Context, state string) error {
	claims := StateClaims{
		State: state,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(o.key)

	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return err
	}
	ctx.SetCookie(o.stateCookieName, tokenStr, 600, "/oauth/wechat/callback", "",
		false, true)
	return nil
}
