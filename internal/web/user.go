package web

import (
	"errors"
	"fmt"
	"github.com/gin-contrib/sessions"
	"net/http"
	"webook/internal/domain"
	"webook/internal/service"

	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// UserHandler 定义用户相关路由
type UserHandler struct {
	svc            *service.UserService
	regexpEmail    *regexp.Regexp
	regexpPassword *regexp.Regexp
}

// NewUserHandler 新建一个UserHandler 包含email 和 password 的正则预编译
func NewUserHandler(svc *service.UserService) *UserHandler {
	// 校验正则表达式 是否写错
	const (
		emailRegex    = `^\w+(-+.\w+)*@\w+(-.\w+)*.\w+(-.\w+)*$`
		passwordRegex = `^(?=.*\d)(?=.*[a-z])(?=.*[A-Z])(?=.*[^a-zA-Z\d]).{8,20}$`
	)

	// 校验正则、邮箱
	regexEmail := regexp.MustCompile(emailRegex, 0)
	// 校验正则、密码
	regexPassword := regexp.MustCompile(passwordRegex, regexp.None)

	return &UserHandler{
		svc:            svc,
		regexpEmail:    regexEmail,
		regexpPassword: regexPassword,
	}
}

// RegisterRoutes 注册用户相关路由
func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	// todo
	server.POST("/users/signup", u.SignUp)
	//server.POST("/users/login", u.Login)
	server.POST("/users/login", u.LoginJWT)
	server.POST("/users/edit", u.Edit)
	server.GET("/users/profile", u.Profile)
}

func (u *UserHandler) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	var req SignUpReq
	// 根据Content-Type 解析数据到 req中，
	// 错误自动写回 400
	// 注意要传地址
	if err := ctx.Bind(&req); err != nil {
		return
	}

	if isMatch, err := u.regexpEmail.MatchString(req.Email); err != nil {
		// todo 记录日志
		ctx.String(http.StatusInternalServerError, "系统错误")
		return
	} else if !isMatch {
		ctx.String(http.StatusBadRequest, "无效邮箱")
		return
	}

	if isMatch, err := u.regexpPassword.MatchString(req.Password); err != nil {
		// todo 记录日志
		ctx.String(http.StatusInternalServerError, "系统错误")
		return
	} else if !isMatch {
		ctx.String(http.StatusBadRequest, "无效密码")
		return
	}

	// 密码确认 校验
	if req.ConfirmPassword != req.Password {
		ctx.String(http.StatusBadRequest, "密码不一致")
		return
	}

	// 调用service
	err := u.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})

	if err == service.ErrUserDuplicateEmail {
		ctx.String(http.StatusInternalServerError, "该邮箱已被注册")
		return
	}
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "注册成功！！！",
	})

}

type UserClaims struct {
	jwt.RegisteredClaims
	UserID    int64
	UserAgent string
}

func (u *UserHandler) LoginJWT(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq

	// 解析JSON数据
	if err := ctx.Bind(&req); err != nil {
		return
	}
	// 身份校验
	user, err := u.svc.Login(ctx, req.Email, req.Password)
	if errors.Is(err, service.ErrInvalidUserOrPassword) {
		ctx.String(http.StatusOK, "邮箱或密码错误")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	// 2.生成JWT
	claims := UserClaims{
		UserID:    user.ID,
		UserAgent: ctx.Request.UserAgent(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString([]byte("Hbzhtd0211"))

	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return
	}
	ctx.Header("x-jwt-token", tokenStr)
	fmt.Println(user)

	ctx.String(http.StatusOK, "登录成功")
	return
}

func (u *UserHandler) Login(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq

	// 解析JSON数据
	if err := ctx.Bind(&req); err != nil {
		return
	}
	// 身份校验
	user, err := u.svc.Login(ctx, req.Email, req.Password)
	if errors.Is(err, service.ErrInvalidUserOrPassword) {
		ctx.String(http.StatusOK, "邮箱或密码错误")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	// 2.设置session
	session := sessions.Default(ctx)
	session.Set("userID", user.ID) // 把userID 存入session
	session.Options(sessions.Options{
		MaxAge: 60 * 30, //30min
	})
	err = session.Save()
	// session 保存失败
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	ctx.String(http.StatusOK, "登录成功")
	return
}

func (u *UserHandler) Logout(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Options(sessions.Options{
		MaxAge: -1,
	})
	err := session.Save()
	if err != nil {
		return
	}
	ctx.String(http.StatusOK, "退出登录成功...")
}

func (u *UserHandler) Edit(ctx *gin.Context) {

}

func (u *UserHandler) Profile(ctx *gin.Context) {
	ctx.String(200, "用户信息...")
}
