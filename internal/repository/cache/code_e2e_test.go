package cache

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRedisCodeCache_Set_e2e(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	testCases := []struct {
		name string

		before func(t *testing.T)
		after  func(t *testing.T)

		ctx   context.Context
		biz   string
		phone string
		code  string

		wantErr error
	}{
		{
			name: "设置成功",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:login:15212345678"

				dur, err := rdb.TTL(ctx, key).Result()
				assert.NoError(t, err)
				assert.True(t, dur > time.Minute*9+time.Second+50)
				code, err := rdb.Get(ctx, key).Result()
				assert.NoError(t, err)
				assert.Equal(t, "123456", code)
			},
			ctx:     context.Background(),
			biz:     "login",
			phone:   "15212345678",
			code:    "123456",
			wantErr: nil,
		},
		{
			name: "发送太频繁",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:login:15212345678"
				err := rdb.Set(ctx, key, "654321", time.Minute*9+time.Second*50).Err()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:login:15212345678"

				dur, err := rdb.TTL(ctx, key).Result()
				assert.NoError(t, err)
				assert.True(t, dur > time.Minute*9+time.Second+50)
				code, err := rdb.Get(ctx, key).Result()
				assert.NoError(t, err)
				assert.Equal(t, "654321", code)
			},
			ctx:     context.Background(),
			biz:     "login",
			phone:   "15212345678",
			code:    "123456",
			wantErr: ErrCodeSendTooMany,
		},
		{
			name: "系统错误",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:login:15212345678"
				err := rdb.Set(ctx, key, "654321", 0).Err()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:login:15212345678"
				code, err := rdb.GetDel(ctx, key).Result()
				assert.NoError(t, err)
				assert.Equal(t, "654321", code)
			},
			ctx:     context.Background(),
			biz:     "login",
			phone:   "15212345678",
			code:    "123456",
			wantErr: errors.New("验证码存在，但是没有过期时间"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			defer tc.after(t)
			c := NewCodeCache(rdb)
			err := c.Set(tc.ctx, tc.biz, tc.phone, tc.code)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
