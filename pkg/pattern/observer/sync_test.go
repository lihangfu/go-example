package observer

import (
	"context"
	"testing"
)

func Test_syncEventBus(t *testing.T) {
	observerA := NewBaseObserver("a")
	observerB := NewBaseObserver("b")
	observerC := NewBaseObserver("c")
	observerD := NewBaseObserver("d")

	sbus := NewSyncEventBus()
	topic := "order_finish"
	sbus.Subscribe(topic, observerA)
	sbus.Subscribe(topic, observerB)
	sbus.Subscribe(topic, observerC)
	sbus.Subscribe(topic, observerD)

	sbus.Publish(context.Background(), &Event{
		Topic: topic,
		Val:   "order_id: xxx",
	})
}
