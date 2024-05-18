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

func TestTimeoutFailOverSMSService_Send(t *testing.T) {
	testCases := []struct {
		name      string
		mock      func(ctrl *gomock.Controller) []sms.Service
		idx       int32
		cnt       int32
		threshold int32
		wantErr   error
		wantIdx   int32
		wantCnt   int32
	}{
		{
			name: "没有触发切换",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc := sms_mocksvc.NewMockService(ctrl)
				svc.EXPECT().
					Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
				return []sms.Service{svc}
			},
			idx:       0,
			cnt:       0,
			threshold: 3,
			wantErr:   nil,
			wantCnt:   0,
			wantIdx:   0,
		},
		{
			name: "触发切换，成功",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := sms_mocksvc.NewMockService(ctrl)
				svc1 := sms_mocksvc.NewMockService(ctrl)
				svc1.EXPECT().Send(gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any()).Return(nil)
				return []sms.Service{svc0, svc1}
			},
			idx:       0,
			cnt:       15,
			threshold: 15,
			// 触发了切换
			wantIdx: 1,
			wantCnt: 0,
			wantErr: nil,
		},
		{
			name: "切换失败",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc := sms_mocksvc.NewMockService(ctrl)
				svc1 := sms_mocksvc.NewMockService(ctrl)
				svc.EXPECT().
					Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("send err"))
				return []sms.Service{svc, svc1}
			},
			idx:       1,
			cnt:       3,
			threshold: 3,
			wantErr:   errors.New("send err"),
			wantIdx:   0,
			wantCnt:   0,
		},
		{
			name: "触发切换，发送超时",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc := sms_mocksvc.NewMockService(ctrl)
				svc1 := sms_mocksvc.NewMockService(ctrl)
				svc1.EXPECT().
					Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(context.DeadlineExceeded)
				return []sms.Service{svc, svc1}
			},
			idx:       0,
			cnt:       3,
			threshold: 3,
			wantErr:   context.DeadlineExceeded,
			wantIdx:   1,
			wantCnt:   1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			svc := NewTimeoutFailOverSMSService(tc.mock(ctrl), tc.threshold)
			svc.cnt = tc.cnt
			svc.idx = tc.idx
			err := svc.Send(context.Background(), "123", []string{"123"}, "123")
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantIdx, svc.idx)
			assert.Equal(t, tc.wantCnt, svc.cnt)
		})
	}
}
