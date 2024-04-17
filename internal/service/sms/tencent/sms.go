package sms

import (
	"context"
	"fmt"
	"github.com/ecodeclub/ekit"
	"github.com/ecodeclub/ekit/slice"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type Service struct {
	appID    *string
	signName *string
	client   *sms.Client
}

func NewService(client *sms.Client, appID string, signName string) *Service {
	return &Service{
		client:   client,
		appID:    ekit.ToPtr[string](appID),
		signName: ekit.ToPtr[string](signName),
	}
}

func (s *Service) Send(ctx context.Context, tplID string, args []string, numbers ...string) error {
	fmt.Println(ctx)
	req := sms.NewSendSmsRequest()
	req.SmsSdkAppId = s.appID
	req.SignName = s.signName
	req.TemplateId = ekit.ToPtr[string](tplID)
	req.PhoneNumberSet = s.toStringPtrSlic(numbers)
	req.TemplateParamSet = s.toStringPtrSlic(args)
	resp, err := s.client.SendSms(req)
	if err != nil {
		return err
	}
	for _, status := range resp.Response.SendStatusSet {
		if status.Code == nil || *(status.Code) != "OK" {
			return fmt.Errorf("发送短信失败 %s %s", *status.Code, *status.Message)
		}
	}
	return nil
}

func (s *Service) toStringPtrSlic(src []string) []*string {
	return slice.Map[string, *string](src, func(idx int, src string) *string {
		return &src
	})
}
