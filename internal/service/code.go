package service

import (
	"context"
	"fmt"
	"math/rand"
	"webook/internal/repository"
	"webook/internal/service/sms/tencent"
)

type CodeService struct {
	repo   *repository.CodeRepository
	smsSvc sms.Service
}

func (svc *CodeService) Send(ctx context.Context, biz, phone string) error {
	// 生成验证码

	// 存入Redis
	// 发送
}

func (svc *CodeService) Verify(ctx context.Context, biz, phone, inputCode string) error {
	code := svc.generateCode()
	err := svc.repo.Store(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	svc.smsSvc.Send(ctx)
}

func (svc *CodeService) generateCode() string {
	num := rand.Intn(1000000)
	return fmt.Sprintf("%06d", num)
}
