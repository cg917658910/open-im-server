package common

import "sync"

type Broker struct {
	mu    sync.RWMutex
	subs  map[chan string]struct{}
	close chan struct{}
}

func NewBroker() *Broker {
	return &Broker{
		subs:  make(map[chan string]struct{}),
		close: make(chan struct{}),
	}
}

func (b *Broker) Subscribe(buffer int) chan string {
	ch := make(chan string, buffer)
	b.mu.Lock()
	b.subs[ch] = struct{}{}
	b.mu.Unlock()
	return ch
}

func (b *Broker) Unsubscribe(ch chan string) {
	b.mu.Lock()
	if _, ok := b.subs[ch]; ok {
		delete(b.subs, ch)
		close(ch)
	}
	b.mu.Unlock()
}

func (b *Broker) Publish(msg string) {
	b.mu.RLock()
	for ch := range b.subs {
		select {
		case ch <- msg:
		default:
		}
	}
	b.mu.RUnlock()
}

func (b *Broker) Close() {
	b.mu.Lock()
	select {
	case <-b.close:
		b.mu.Unlock()
		return
	default:
	}
	close(b.close)
	for ch := range b.subs {
		close(ch)
		delete(b.subs, ch)
	}
	b.mu.Unlock()
}
