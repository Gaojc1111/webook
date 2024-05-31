package web

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"strings"
	"time"
)

type JWTHandler struct {
	signingMethod jwt.SigningMethod
	access_key    []byte
	refresh_key   []byte
}

type UserClaims struct {
	jwt.RegisteredClaims
	UserID    int64
	UserAgent string
}

type RefreshClaims struct {
	jwt.RegisteredClaims
	UserID    int64
	UserAgent string
}

func NewJWTHandler() JWTHandler {
	return JWTHandler{
		signingMethod: jwt.SigningMethodHS512,
		access_key:    []byte("Hbzhtd0211"),
		refresh_key:   []byte("Gaojc1111"),
	}
}

func (j *JWTHandler) setJWTToken(ctx *gin.Context, userID int64) error {
	if err := j.setAccessJWTToken(ctx, userID); err != nil {
		return err
	}
	if err := j.setRefreshJWTToken(ctx, userID); err != nil {
		return err
	}
	return nil
}

func (j *JWTHandler) setAccessJWTToken(ctx *gin.Context, userID int64) error {
	claims := UserClaims{
		UserID:    userID,
		UserAgent: ctx.Request.UserAgent(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	token := jwt.NewWithClaims(j.signingMethod, claims)
	tokenStr, err := token.SignedString(j.access_key)

	if err != nil {
		return err
	}
	ctx.Header("x-jwt-token", tokenStr)
	return nil
}

func (j *JWTHandler) setRefreshJWTToken(ctx *gin.Context, userID int64) error {
	claims := RefreshClaims{
		UserID:    userID,
		UserAgent: ctx.Request.UserAgent(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
	}
	token := jwt.NewWithClaims(j.signingMethod, claims)
	tokenStr, err := token.SignedString(j.refresh_key)
	if err != nil {
		return err
	}
	ctx.Header("x-refresh-token", tokenStr)
	return nil
}

func ParseToken(ctx *gin.Context) string {
	tokenStr := ctx.GetHeader("Authorization")
	if tokenStr == "" {
		return tokenStr
	}
	segs := strings.Split(tokenStr, " ") // Bearer token...
	if len(segs) != 2 {
		return tokenStr
	}
	return segs[1]
}
