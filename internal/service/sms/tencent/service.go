package tencent

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type Service struct {
	client    *sms.Client
	appID     *string
	signature *string
}

func (s *Service) Send(ctx context.Context, tplID string, args []string, numbers ...string) error {
	request := sms.NewSendSmsRequest()
	request.SmsSdkAppId = common.StringPtr("1400787878")
	request.SignName = common.StringPtr("腾讯云")
	request.TemplateId = common.StringPtr(tplID)
	request.TemplateParamSet = common.StringPtrs(args)
	request.PhoneNumberSet = common.StringPtrs(numbers)
	response, err := s.client.SendSms(request)

	if err != nil {
		fmt.Printf("%s", err.Error())
		return err
	}
	for _, status := range response.Response.SendStatusSet {
		if *status.Code != "Ok" {
			return fmt.Errorf("send sms failed, code=%s, message=%s", *status.Code, *status.Message)
		}
	}
	b, _ := json.Marshal(response.Response)
	// 打印返回的json字符串
	fmt.Printf("%s", b)
	return nil
}

func NewService(client *sms.Client, appID string, signName string) *Service {
	return &Service{
		client:    client,
		appID:     &appID,
		signature: &signName,
	}
}
