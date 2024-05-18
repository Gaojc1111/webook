package failover

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"webook/internal/service/sms"
	"webook/internal/service/sms/sms_mocksvc"
)

func TestFailOverSMSService_Send(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) []sms.Service
		ctx     context.Context
		tplID   string
		args    []string
		numbers []string
		wantErr error
	}{
		{
			name: "一次成功",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				s1 := sms_mocksvc.NewMockService(ctrl)
				s1.EXPECT().
					Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
				return []sms.Service{s1}
			},
			wantErr: nil,
		},
		{
			name: "两次成功",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				s1 := sms_mocksvc.NewMockService(ctrl)
				s1.EXPECT().
					Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("err"))
				s2 := sms_mocksvc.NewMockService(ctrl)
				s2.EXPECT().
					Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
				return []sms.Service{s1, s2}
			},
			wantErr: nil,
		},
		{
			name: "全部失败",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				s1 := sms_mocksvc.NewMockService(ctrl)
				s1.EXPECT().
					Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("err"))
				s2 := sms_mocksvc.NewMockService(ctrl)
				s2.EXPECT().
					Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("err"))
				return []sms.Service{s1, s2}
			},
			wantErr: errors.New("所有服务商都发送失败"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			sms_arr := tc.mock(ctrl)

			svc := NewFailOverSMSService(sms_arr)
			err := svc.Send(tc.ctx, tc.tplID, tc.args, tc.numbers...)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
