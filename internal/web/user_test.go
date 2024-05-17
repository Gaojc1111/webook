package web

import (
	"bytes"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
	"webook/internal/domain"
	"webook/internal/service"
	mocksvc "webook/internal/service/mock"
)

// 测试UserHandler注册路由
func TestUserHandler_SignUp(t *testing.T) {
	testCases := []struct {
		name string
		// mock 依赖
		mock func(ctrl *gomock.Controller) (service.UserService, service.CodeService)
		// 预期请求
		reqBuild func(t *testing.T) *http.Request
		// 预期响应
		wantCode int
		wantBody string
	}{
		{
			name: "注册成功",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := mocksvc.NewMockUserService(ctrl)
				codeSvc := mocksvc.NewMockCodeService(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "111@qq.com",
					Password: "QQqq11!!",
				}).Return(nil)
				return userSvc, codeSvc
			},
			reqBuild: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte(`{
"email": "111@qq.com",
"password":"QQqq11!!",
"confirmPassword":"QQqq11!!"
}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},
			wantCode: 200,
			wantBody: "注册成功",
		},
		{
			name: "Bind 失败",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := mocksvc.NewMockUserService(ctrl)
				codeSvc := mocksvc.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuild: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte(`{
"emaill": "111@qq.com",
"password":"QQqq11!"
"confirmPassword":"QQqq11!!，"
}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},
			wantCode: 400,
		},
		{
			name: "无效邮箱",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				return nil, nil
			},
			reqBuild: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte(`{
"email": "11@",
"password":"QQqq11!!",
"confirmPassword":"QQqq11!!"
}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},
			wantCode: 400,
			wantBody: "无效邮箱",
		},
		{
			name: "无效密码",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				return nil, nil
			},
			reqBuild: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte(`{
"email": "111@qq.com",
"password":"1111",
"confirmPassword":"1111"
}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},
			wantCode: 400,
			wantBody: "无效密码",
		},
		{
			name: "密码不一致",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				return nil, nil
			},
			reqBuild: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte(`{
"email": "111@qq.com",
"password":"QQqq11!!",
"confirmPassword":"QQqq11!"
}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},
			wantCode: 400,
			wantBody: "密码不一致",
		},
		{
			name: "该邮箱已被注册",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := mocksvc.NewMockUserService(ctrl)
				codeSvc := mocksvc.NewMockCodeService(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "111@qq.com",
					Password: "QQqq11!!",
				}).Return(service.ErrUserDuplicated)
				return userSvc, codeSvc
			},
			reqBuild: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte(`{
"email": "111@qq.com",
"password":"QQqq11!!",
"confirmPassword":"QQqq11!!"
}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},
			wantCode: 400,
			wantBody: "该邮箱已被注册",
		},
		{
			name: "系统错误",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := mocksvc.NewMockUserService(ctrl)
				codeSvc := mocksvc.NewMockCodeService(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "111@qq.com",
					Password: "QQqq11!!",
				}).Return(errors.New("DB err"))
				return userSvc, codeSvc
			},
			reqBuild: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte(`{
"email": "111@qq.com",
"password":"QQqq11!!",
"confirmPassword":"QQqq11!!"
}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},
			wantCode: 500,
			wantBody: "系统错误",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			//mock需要的service
			userSvc, codeSvc := tc.mock(ctrl)
			hdl := NewUserHandler(userSvc, codeSvc)

			// 构造server & 注册路由
			server := gin.Default()
			hdl.RegisterRoutes(server)

			req := tc.reqBuild(t)
			recorder := httptest.NewRecorder()

			// 执行测试
			server.ServeHTTP(recorder, req)

			// 断言结果
			assert.Equal(t, tc.wantCode, recorder.Code)
			assert.Equal(t, tc.wantBody, recorder.Body.String())
		})
	}
}

func TestMock(t *testing.T) {
	// mock使用
	// 初始化控制器
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// 模拟实例
	userServer := mocksvc.NewMockUserService(ctrl)
	// 模拟调用， 不管第一个参数
	userServer.EXPECT().SignUp(gomock.Any(), domain.User{
		ID:    1,
		Email: "666@qq.com",
	}).Return(errors.New("bad param"))

	err := userServer.SignUp(context.Background(), domain.User{
		ID:    1,
		Email: "666@qq.com",
	})
	t.Log(err)
}
