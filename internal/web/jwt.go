package web

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

type JWTHandler struct {
}

type UserClaims struct {
	jwt.RegisteredClaims
	UserID    int64
	UserAgent string
}

func NewJWTHandler() *JWTHandler {
	return &JWTHandler{}
}

func (j *JWTHandler) setJWTToken(ctx *gin.Context, userID int64) {
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
