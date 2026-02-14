package common

import (
	"testing"
	"time"
)

func TestBrokerPublishSubscribe(t *testing.T) {
	b := NewBroker()
	defer b.Close()

	sub := b.Subscribe(1)
	defer b.Unsubscribe(sub)

	b.Publish("x")
	select {
	case msg := <-sub:
		if msg != "x" {
			t.Fatalf("unexpected msg: %s", msg)
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting message")
	}
}
