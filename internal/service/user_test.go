package service

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"webook/internal/domain"
	"webook/internal/repository"
	mocksvc "webook/internal/repository/mock"
)

func Test_userService_Login(t *testing.T) {
	testCases := []struct {
		name     string
		mock     func(controller *gomock.Controller) repository.UserRepository
		ctx      context.Context
		email    string
		password string
		wantUser domain.User
		wantErr  error
	}{
		{
			name: "登陆成功",
			ctx:  context.Background(),
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				userRepository := mocksvc.NewMockUserRepository(ctrl)
				userRepository.
					EXPECT().FindByEmail(gomock.Any(), "666@qq.com").
					Return(domain.User{
						Email:    "666@qq.com",
						Password: "$2a$10$EHqoKRCV1mAPyeUebxNUeeOK2lAGvpsxT1pUZFgvw9TuKA9EVNLvS",
					},
						nil,
					)
				return userRepository
			},
			email:    "666@qq.com",
			password: "QQqq11!!",

			wantUser: domain.User{
				Email:    "666@qq.com",
				Password: "$2a$10$EHqoKRCV1mAPyeUebxNUeeOK2lAGvpsxT1pUZFgvw9TuKA9EVNLvS",
			},
			wantErr: nil,
		},
		{
			name: "查询不到该用户",
			ctx:  context.Background(),
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				userRepository := mocksvc.NewMockUserRepository(ctrl)
				userRepository.
					EXPECT().FindByEmail(gomock.Any(), "111@qq.com").
					Return(domain.User{}, repository.ErrUserNotFound)
				return userRepository
			},
			email:    "111@qq.com",
			password: "QQqq11!!",
			wantUser: domain.User{},
			wantErr:  ErrInvalidUserOrPassword,
		},
		{
			name: "系统错误",
			ctx:  context.Background(),
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				userRepository := mocksvc.NewMockUserRepository(ctrl)
				userRepository.
					EXPECT().FindByEmail(gomock.Any(), "111@qq.com").
					Return(domain.User{}, errors.New("db err"))
				return userRepository
			},
			email:    "111@qq.com",
			password: "QQqq11!!",
			wantUser: domain.User{},
			wantErr:  errors.New("db err"),
		},
		{
			name: "密码错误",
			ctx:  context.Background(),
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				userRepository := mocksvc.NewMockUserRepository(ctrl)
				userRepository.
					EXPECT().FindByEmail(gomock.Any(), "666@qq.com").
					Return(domain.User{
						Email:    "666@qq.com",
						Password: "$2a$10$EHqoKRCV1mAPyeUebxNUeeOK2lAGvpsxT1pUZFgvw9TuKA9EVNLvS",
					}, nil)
				return userRepository
			},
			email:    "666@qq.com",
			password: "QQqq11!",
			wantErr:  ErrInvalidUserOrPassword,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepository := tc.mock(ctrl)
			userSvc := NewUserService(userRepository)

			u, err := userSvc.Login(tc.ctx, tc.email, tc.password)

			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, u)
		})
	}
}
