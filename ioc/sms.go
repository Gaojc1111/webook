package ioc

import (
	"webook/internal/service/sms"
	"webook/internal/service/sms/localsms"
)

func InitSMSService() sms.Service {
	return localsms.NewService()
}
