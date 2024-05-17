package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"webook/internal/repository/cache/redis_mock"
)

func TestRedisCodeCache_Set(t *testing.T) {
	keyFunc := func(biz, phone string) string {
		return fmt.Sprintf("phone_code:%s:%s", biz, phone)
	}
	testCases := []struct {
		keyFunc func(ctrl *gomock.Controller) redis.Cmdable
		name    string
		mock    func(ctrl *gomock.Controller) redis.Cmdable
		ctx     context.Context
		biz     string
		phone   string
		code    string
		wantErr error
	}{
		{
			name: "缓存成功",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				res := redis_mock.NewMockCmdable(ctrl)
				cmd := redis.NewCmd(context.Background())
				cmd.SetErr(nil)
				cmd.SetVal(int64(0))
				res.EXPECT().
					Eval(gomock.Any(), luaSetCode, []string{(keyFunc("test", "12312345678"))}, "123").
					Return(cmd)
				return res
			},
			ctx:     context.Background(),
			biz:     "test",
			phone:   "12312345678",
			code:    "123",
			wantErr: nil,
		},
		{
			name: "redis err",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				res := redis_mock.NewMockCmdable(ctrl)
				cmd := redis.NewCmd(context.Background())
				err := errors.New("redis err")
				cmd.SetErr(err)
				cmd.SetVal(int64(0))
				res.EXPECT().
					Eval(gomock.Any(), luaSetCode, []string{(keyFunc("test", "12312345678"))}, "123").
					Return(cmd)
				return res
			},
			ctx:     context.Background(),
			biz:     "test",
			phone:   "12312345678",
			code:    "123",
			wantErr: errors.New("redis err"),
		},
		{
			name: "验证码存在，但是没有过期时间",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				res := redis_mock.NewMockCmdable(ctrl)
				cmd := redis.NewCmd(context.Background())
				err := errors.New("验证码存在，但是没有过期时间")
				cmd.SetErr(err)
				cmd.SetVal(int64(-2))
				res.EXPECT().
					Eval(gomock.Any(), luaSetCode, []string{(keyFunc("test", "12312345678"))}, "123").
					Return(cmd)
				return res
			},
			ctx:     context.Background(),
			biz:     "test",
			phone:   "12312345678",
			code:    "123",
			wantErr: errors.New("验证码存在，但是没有过期时间"),
		},
		{
			name: "验证码发送太频繁",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				res := redis_mock.NewMockCmdable(ctrl)
				cmd := redis.NewCmd(context.Background())
				err := ErrCodeSendTooMany
				cmd.SetErr(err)
				cmd.SetVal(int64(-1))
				res.EXPECT().
					Eval(gomock.Any(), luaSetCode, []string{(keyFunc("test", "12312345678"))}, "123").
					Return(cmd)
				return res
			},
			ctx:     context.Background(),
			biz:     "test",
			phone:   "12312345678",
			code:    "123",
			wantErr: ErrCodeSendTooMany,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			codeCache := NewCodeCache(tc.mock(ctrl))
			err := codeCache.Set(tc.ctx, tc.biz, tc.phone, tc.code)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
