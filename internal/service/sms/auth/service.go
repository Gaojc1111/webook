package auth

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"webook/internal/service/sms"
)

type SMSService struct {
	svc sms.Service
	key []byte
}

func NewSMSService(sms sms.Service, key []byte) *SMSService {
	return &SMSService{
		svc: sms,
		key: key,
	}
}

func (s SMSService) Send(ctx context.Context, tplToken string, args []string, numbers ...string) error {
	var claims SMSClaims
	_, err := jwt.ParseWithClaims(tplToken, &claims, func(token *jwt.Token) (interface{}, error) {
		return s.key, nil
	})
	if err != nil {
		return err
	}
	return s.svc.Send(ctx, claims.tpl, args, numbers...)
}

type SMSClaims struct {
	tpl string
	jwt.RegisteredClaims
}
