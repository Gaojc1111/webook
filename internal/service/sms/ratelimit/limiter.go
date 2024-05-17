package ratelimit

import (
	"context"
	"errors"
	"webook/internal/service/sms"
	"webook/pkg/limiter"
)

var ErrRateLimit = errors.New("触发限流")

type RateLimitSMSService struct {
	svc     sms.Service
	limiter limiter.Limiter
	key     string
}

func (r RateLimitSMSService) Send(ctx context.Context, tplID string, args []string, numbers ...string) error {
	limited, err := r.limiter.Limit(ctx, r.key)
	if err != nil {
		return err
	}
	if limited {
		return ErrRateLimit
	}
	return r.svc.Send(ctx, tplID, args, numbers...)
}

func NewRateLimitSMSService(svc sms.Service, limiter limiter.Limiter) *RateLimitSMSService {
	return &RateLimitSMSService{
		svc:     svc,
		limiter: limiter,
		key:     "sms-limiter",
	}
}
