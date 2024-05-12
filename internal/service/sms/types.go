package sms

import "context"

// Service 短信服务的抽象
type Service interface {
	Send(ctx context.Context, tplID string, args []string, numbers ...string) error
}
