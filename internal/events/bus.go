package events

import (
    "reflect"
    "sync"
)

type Event interface{}

type HandlerFunc func(Event)

type Bus struct {
    subs map[string][]HandlerFunc
    mu   sync.RWMutex
}

func NewBus() *Bus {
    return &Bus{
        subs: make(map[string][]HandlerFunc),
    }
}

func (b *Bus) Subscribe(eventName string, h HandlerFunc) {
    b.mu.Lock()
    defer b.mu.Unlock()
    b.subs[eventName] = append(b.subs[eventName], h)
}

func (b *Bus) Publish(event Event) {
    name := reflect.TypeOf(event).Name()
    b.mu.RLock()
    handlers := b.subs[name]
    b.mu.RUnlock()
    for _, h := range handlers {
        go h(event) // async handler
    }
}
