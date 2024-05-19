package ioc

import (
	"webook/internal/service/oauth2/wechat"
)

func InitWechatService() wechat.Service {
	//appID, ok := os.LookupEnv("WECHAT_APP_ID")
	//if !ok {
	//	panic("WECHAT_APP_ID not found")
	//}
	//appSecret, ok := os.LookupEnv("WECHAT_APP_SECRET")
	//if !ok {
	//	panic("WECHAT_APP_SECRET not found")
	//}
	// todo 微信开发者appID 有时间再搞吧...
	appID := "123"
	appSecret := "123"
	return wechat.NewService(appID, appSecret)
}
