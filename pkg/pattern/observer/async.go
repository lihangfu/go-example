package observer

import (
	"context"
	"fmt"
)

type observerWithErr struct {
	o   Observer
	err error
}

type AsyncEventBus struct {
	BaseEventBus
	errC chan *observerWithErr
	ctx  context.Context
	stop context.CancelFunc
}

func NewAsyncEventBus() *AsyncEventBus {
	aBus := AsyncEventBus{
		BaseEventBus: NewBaseEventBus(),
	}
	aBus.ctx, aBus.stop = context.WithCancel(context.Background())
	// 处理处理错误的异步守护协程
	go aBus.handleErr()
	return &aBus
}

func (a *AsyncEventBus) Stop() {
	a.stop()
}

func (a *AsyncEventBus) Publish(ctx context.Context, e *Event) {
	a.mux.RLock()
	subscribers := a.observers[e.Topic]
	a.mux.RUnlock()
	for subscriber := range subscribers {
		// shadow
		subscriber := subscriber
		go func() {
			if err := subscriber.OnChange(ctx, e); err != nil {
				select {
				case <-a.ctx.Done():
				case a.errC <- &observerWithErr{
					o:   subscriber,
					err: err,
				}:
				}
			}
		}()
	}
}

func (a *AsyncEventBus) handleErr() {
	for {
		select {
		case <-a.ctx.Done():
			return
		case resp := <-a.errC:
			// 处理 publish 失败的 observer
			fmt.Printf("observer: %v, err: %v", resp.o, resp.err)
		}
	}
}
