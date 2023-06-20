package observer

import (
	"context"
	"fmt"
	"sync"
)

// Event 事物的变更事件. 其中 Topic 标识了事物的身份以及变更的类型，Val 是变更详情
type Event struct {
	Topic string
	Val   interface{}
}

// Observer 观察者. 指的是关注事物动态的角色
type Observer interface {
	OnChange(ctx context.Context, e *Event) error
}

// EventBus 事件总线. 位于观察者与事物之间承上启下的代理层. 负责维护管理观察者，并且在事物发生变更时，将情况同步给每个观察者.
type EventBus interface {
	// Subscribe 针对于观察者而言，需要向 EventBus 完成注册操作，注册时需要声明自己关心的变更事件类型，不再需要直接和事物打交道
	Subscribe(topic string, o Observer)
	// Unsubscribe 取消订阅
	Unsubscribe(topic string, o Observer)
	// Publish 针对于事物而言，在其发生变更时，只需要将变更情况向 EventBus 统一汇报即可，不再需要和每个观察者直接交互
	Publish(ctx context.Context, e *Event)
}

var _ Observer = (*BaseObserver)(nil)

type BaseObserver struct {
	name string
}

func NewBaseObserver(name string) *BaseObserver {
	return &BaseObserver{
		name: name,
	}
}

func (b *BaseObserver) OnChange(ctx context.Context, e *Event) error {
	_ = ctx
	fmt.Printf("observer: %s, event key: %s, event val: %v", b.name, e.Topic, e.Val)
	return nil
}

var _ EventBus = (*BaseEventBus)(nil)

type BaseEventBus struct {
	mux       sync.RWMutex
	observers map[string]map[Observer]struct{}
}

func NewBaseEventBus() BaseEventBus {
	return BaseEventBus{
		observers: make(map[string]map[Observer]struct{}),
	}
}

func (b *BaseEventBus) Subscribe(topic string, o Observer) {
	b.mux.Lock()
	defer b.mux.Unlock()
	_, ok := b.observers[topic]
	if !ok {
		b.observers[topic] = make(map[Observer]struct{})
	}
	b.observers[topic][o] = struct{}{}
}

func (b *BaseEventBus) Unsubscribe(topic string, o Observer) {
	b.mux.Lock()
	defer b.mux.Unlock()
	delete(b.observers[topic], o)
}

func (b *BaseEventBus) Publish(ctx context.Context, e *Event) {
	_, _ = ctx, e
	panic("implement me")
}
