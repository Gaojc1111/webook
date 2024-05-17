package ratelimit

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"webook/internal/service/sms"
	"webook/internal/service/sms/sms_mocksvc"
	"webook/pkg/limiter"
	limiter_mocksvc "webook/pkg/limiter/mock"
)

func TestRateLimitSMSService_Send(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) (svc sms.Service, limiter limiter.Limiter)
		wantErr error
	}{
		{
			name: "不限流",
			mock: func(ctrl *gomock.Controller) (sms.Service, limiter.Limiter) {
				svc := sms_mocksvc.NewMockService(ctrl)
				l := limiter_mocksvc.NewMockLimiter(ctrl)
				l.EXPECT().Limit(gomock.Any(), gomock.Any()).Return(false, nil)
				svc.EXPECT().
					Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
				return svc, l
			},
			wantErr: nil,
		},
		{
			name: "限流",
			mock: func(ctrl *gomock.Controller) (sms.Service, limiter.Limiter) {
				svc := sms_mocksvc.NewMockService(ctrl)
				l := limiter_mocksvc.NewMockLimiter(ctrl)
				l.EXPECT().Limit(gomock.Any(), gomock.Any()).Return(true, nil)
				return svc, l
			},
			wantErr: ErrRateLimit,
		},
		{
			name: "限流器错误",
			mock: func(ctrl *gomock.Controller) (sms.Service, limiter.Limiter) {
				svc := sms_mocksvc.NewMockService(ctrl)
				l := limiter_mocksvc.NewMockLimiter(ctrl)
				l.EXPECT().Limit(gomock.Any(), gomock.Any()).Return(false, errors.New("limiter err"))
				return svc, l
			},
			wantErr: errors.New("limiter err"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			svc, limiter := tc.mock(ctrl)
			lss := NewRateLimitSMSService(svc, limiter)
			err := lss.Send(context.Background(), "abc", []string{"123"}, "123")
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
