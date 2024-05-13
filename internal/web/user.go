package web

import (
	"errors"
	"github.com/gin-contrib/sessions"
	"net/http"
	"time"
	"webook/internal/domain"
	"webook/internal/service"

	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	emailRegex    = `^\w+(-+.\w+)*@\w+(-.\w+)*.\w+(-.\w+)*$`
	passwordRegex = `^(?=.*\d)(?=.*[a-z])(?=.*[A-Z])(?=.*[^a-zA-Z\d]).{8,20}$`
	bizLogin      = "Login"
)

// UserHandler 定义用户相关路由
type UserHandler struct {
	svc            *service.UserService
	codeSvc        *service.CodeService
	regexpEmail    *regexp.Regexp
	regexpPassword *regexp.Regexp
}

// NewUserHandler 新建一个UserHandler 包含email 和 password 的正则预编译
func NewUserHandler(svc *service.UserService, codeSvc *service.CodeService) *UserHandler {
	// 正则格式校验
	regexEmail := regexp.MustCompile(emailRegex, 0)
	regexPassword := regexp.MustCompile(passwordRegex, regexp.None)

	return &UserHandler{
		svc:            svc,
		codeSvc:        codeSvc,
		regexpEmail:    regexEmail,
		regexpPassword: regexPassword,
	}
}

// RegisterRoutes 注册用户相关路由
func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	// todo
	ug := server.Group("/users")
	{
		ug.POST("/signup", u.SignUp)
		ug.POST("/login", u.LoginJWT)
		ug.POST("/logout", u.Logout)
		ug.POST("/edit", u.Edit)
		ug.GET("/profile", u.Profile)
	}
	{
		ug.POST("/login_sms/code/send", u.SendSmsCode) // 获取验证码
		ug.POST("/login_sms", u.LoginBySMS)            // 校验验证码
	}
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

	if err == service.ErrUserDuplicated {
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
	switch err {
	case nil:
		u.setJWTToken(ctx, user.ID)
		ctx.String(http.StatusOK, "登录成功")
	case service.ErrInvalidUserOrPassword:
		ctx.String(http.StatusOK, "邮箱或密码错误")
	default:
		ctx.String(http.StatusOK, "系统错误")
	}
}

func (u *UserHandler) setJWTToken(ctx *gin.Context, userID int64) {
	claims := UserClaims{
		UserID:    userID,
		UserAgent: ctx.Request.UserAgent(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString([]byte("Hbzhtd0211"))

	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
	}
	ctx.Header("x-jwt-token", tokenStr)
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
	ctx.String(200, "profile")
}

func (h *UserHandler) LoginBySMS(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}

	ok, err := h.codeSvc.Verify(ctx, bizLogin, req.Phone, req.Code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统异常",
		})
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "验证码不对，请重新输入",
		})
		return
	}
	u, err := h.svc.FindOrCreate(ctx, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	h.setJWTToken(ctx, u.ID)
	ctx.JSON(http.StatusOK, Result{
		Code: 200,
		Msg:  "登录成功",
	})
}

func (h *UserHandler) SendSmsCode(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	// 你这边可以校验 Req
	if req.Phone == "" {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "请输入手机号码",
		})
		return
	}
	err := h.codeSvc.Send(ctx, bizLogin, req.Phone)
	switch err {
	case nil:
		ctx.JSON(http.StatusOK, Result{
			Msg: "发送成功",
		})
	case service.ErrCodeSendTooMany:
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "短信发送太频繁，请稍后再试",
		})
	default:
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		// todo：补日志
	}
}
