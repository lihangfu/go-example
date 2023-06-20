package observer

import (
	"context"
	"fmt"
)

type SyncEventBus struct {
	BaseEventBus
}

func NewSyncEventBus() *SyncEventBus {
	return &SyncEventBus{
		BaseEventBus: NewBaseEventBus(),
	}
}

func (s *SyncEventBus) Publish(ctx context.Context, e *Event) {
	s.mux.RLock()
	subscribers := s.observers[e.Topic]
	s.mux.RUnlock()
	errs := make(map[Observer]error)
	for subscriber := range subscribers {
		if err := subscriber.OnChange(ctx, e); err != nil {
			errs[subscriber] = err
		}
	}
	s.handleErr(ctx, errs)
}

func (s *SyncEventBus) handleErr(ctx context.Context, errs map[Observer]error) {
	_ = ctx
	for o, err := range errs {
		// 处理 publish 失败的 observer
		fmt.Printf("observer: %v, err: %v", o, err)
	}
}
