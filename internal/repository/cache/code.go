package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
)

// 编译器会在编译时，将set_code.lua的代码放进 luaSetCode 变量里
//
//go:embed lua/set_code.lua
var luaSetCode string

//go:embed lua/verify_code.lua
var luaVerifyCode string

var (
	ErrCodeSendTooMany   = errors.New("发送验证码太频繁")
	ErrCodeVerifyTooMany = errors.New("验证码验证次数太多")
)

type CodeCache interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, code string) (bool, error)
}

type RedisCodeCache struct {
	client redis.Cmdable
}

func NewCodeCache(client redis.Cmdable) CodeCache {
	return &RedisCodeCache{
		client: client,
	}
}

// Set 将验证码存储到缓存中。
// ctx: 上下文，用于控制请求的生命周期。
// biz: 业务标识。
// phone: 手机号码。
// code: 验证码。
// 返回值: 错误信息，如果操作成功，则返回nil；否则返回相应的错误。
func (c *RedisCodeCache) Set(ctx context.Context, biz, phone, code string) error {
	// 使用Lua脚本将验证码设置到Redis缓存中，并返回结果。
	res, err := c.client.Eval(ctx, luaSetCode, []string{c.Key(biz, phone)}, code).Int()
	if err != nil {
		return err // 如果有错误发生，则直接返回错误。
	}
	switch res {
	case -2:
		return errors.New("验证码存在，但是没有过期时间")
	case -1:
		// 发送验证码过于频繁的错误。
		return ErrCodeSendTooMany
	default:
		return nil
	}
}

func (c *RedisCodeCache) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	// 获取缓存中的验证码。
	res, err := c.client.Eval(ctx, luaVerifyCode, []string{c.Key(biz, phone)}, code).Int()
	if err != nil {
		return false, err
	}
	switch res {
	case -2:
		return false, nil
	case -1:
		return false, nil
	default:
		return true, nil
	}
}

func (c *RedisCodeCache) Key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}
