package failover

import (
	"context"
	"errors"
	"log"
	"sync/atomic"
	"webook/internal/service/sms"
)

type FailOverSMSService struct {
	svcs []sms.Service
	idx  uint64
}

func NewFailOverSMSService(svcs []sms.Service) *FailOverSMSService {
	return &FailOverSMSService{
		svcs: svcs,
	}
}

// 轮询不同服务商
// 大部分都是在0号服务器
func (f *FailOverSMSService) Send(ctx context.Context, tplID string, args []string, numbers ...string) error {
	for _, svc := range f.svcs {
		err := svc.Send(ctx, tplID, args, numbers...)
		if err == nil {
			return nil
		}
	}
	return errors.New("所有服务商都发送失败")
}

// 优化， 对起始svc进行轮询，避免所有请求都打到0号服务器
func (f *FailOverSMSService) SendV1(ctx context.Context, tplID string, args []string, numbers ...string) error {
	idx := atomic.AddUint64(&f.idx, 1)
	length := uint64(len(f.svcs))
	for i := idx; i < length; i++ {
		err := f.svcs[i%length].Send(ctx, tplID, args, numbers...)
		switch err {
		case nil:
			return nil
		case context.Canceled, context.DeadlineExceeded:
			return err
		}
		log.Println(err)
	}
	return errors.New("所有服务商都发送失败")
}
