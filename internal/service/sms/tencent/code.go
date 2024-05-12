package tencent

import (
	"context"
	"webook/internal/repository"
	"webook/internal/service/sms"
)

type CodeService struct {
	repo *repository.CodeRepository
	sms  sms.Service
}

func (svc *CodeService) Send(ctx context.Context, biz, phone string) error {
	code := svc.GenerateCode()
	err := svc.repo.Set(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	const tplID = "123456"
	return svc.sms.Send(ctx, tplID, []string{code}, phone)
}

func (svc *CodeService) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	ok, err := svc.repo.Verify(ctx, biz, phone, inputCode)
	if err == repository.ErrCodeVerifyTooMany {
		// 返回nil 对业务方屏蔽了错误类型
		return false, nil
	}
	return ok, err
}

func (svc *CodeService) GenerateCode() string {
	return "123456"
}
