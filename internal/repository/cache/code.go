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

var (
	ErrSetCodeSendTooMany = errors.New("发送验证码太频繁")
)

type CodeCache struct {
	client redis.Cmdable
}

func NewCodeCache(client redis.Cmdable) *CodeCache {
	return &CodeCache{
		client: client,
	}
}

func (c *CodeCache) Set(ctx context.Context, biz, phone, code string) error {
	res, err := c.client.Eval(ctx, luaSetCode, []string{c.Key(biz, phone)}, code).Int()
	if err != nil {
		return err
	}
	switch res {
	case 0:
		// 没有问题
		return nil
	case -1:
		return ErrSetCodeSendTooMany
	default:
		return errors.New("系统错误")
	}
}

func (c *CodeCache) Key(biz, phone string) string {
	return fmt.Sprintf("phone:_code:%s:%s", biz, phone)
}
