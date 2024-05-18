package failover

import (
	"context"
	"sync/atomic"
	"webook/internal/service/sms"
)

type TimeoutFailOverSMSService struct {
	svcs      []sms.Service
	idx       int32 // 当前服务索引
	cnt       int32 // 当前连续超时次数
	threshold int32 // 连续超时次数阈值
}

func NewTimeoutFailOverSMSService(svcs []sms.Service, threshod int32) *TimeoutFailOverSMSService {
	return &TimeoutFailOverSMSService{
		svcs:      svcs,
		threshold: threshod,
	}
}

func (t *TimeoutFailOverSMSService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	idx := atomic.LoadInt32(&t.idx)
	cnt := atomic.LoadInt32(&t.cnt)
	// 超过阈值，执行切换
	if cnt >= t.threshold {
		newIdx := (idx + 1) % int32(len(t.svcs))
		if atomic.CompareAndSwapInt32(&t.idx, idx, newIdx) {
			// 重置这个 cnt 计数
			atomic.StoreInt32(&t.cnt, 0)
		}
		idx = newIdx
	}
	svc := t.svcs[idx]
	err := svc.Send(ctx, tplId, args, numbers...)
	switch err {
	case nil:
		// 连续超时，所以不超时的时候要重置到 0
		atomic.StoreInt32(&t.cnt, 0)
		return nil
	case context.DeadlineExceeded:
		atomic.AddInt32(&t.cnt, 1)
	default:
		// 遇到了错误，但是又不是超时错误，这个时候，你要考虑怎么搞
		// 我可以增加，也可以不增加
		// 如果强调一定是超时，那么就不增加
		// 如果是 EOF 之类的错误，你还可以考虑直接切换
	}
	return err
}
