package order

import (
	"github.com/albus-droid/Capstone-Project-Backend/internal/events"
)

type service struct {
    bus *events.Bus
    store map[string]Order // your in-memory store
}

func NewService(bus *events.Bus) Service {
    return &service{
        bus:   bus,
        store: make(map[string]Order),
    }
}

func (s *service) Create(o Order) error {
    if _, exists := s.store[o.ID]; exists {
        return fmt.Errorf("order already exists")
    }
    s.store[o.ID] = o

    // Publish event
    s.bus.Publish(events.OrderPlacedEvent{
        OrderID:   o.ID,
        UserEmail: o.UserEmail,
    })

    return nil
}
