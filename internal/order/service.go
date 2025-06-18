package order

import (
    "errors"
    "fmt"
    "sync"
	"github.com/albus-droid/Capstone-Project-Backend/internal/events"
)

type service struct {
    mu     sync.RWMutex
    store  map[string]Order
    bus    *events.Bus
}

func NewService(bus *events.Bus) Service {
    return &service{
        store: make(map[string]Order),
        bus:   bus,
    }
}

func (s *service) Create(o Order) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    if _, exists := s.store[o.ID]; exists {
        return errors.New("order already exists")
    }

    s.store[o.ID] = o

    // 🔔 Publish async event
    s.bus.Publish(events.OrderPlacedEvent{
        OrderID:   o.ID,
        SellerID:  o.SellerID,
        UserEmail: o.UserEmail,
    })

    fmt.Printf("✅ Order %s created and event published\n", o.ID)
    return nil
}

func (s *service) GetByID(id string) (*Order, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    o, ok := s.store[id]
    if !ok {
        return nil, errors.New("order not found")
    }
    return &o, nil
}

func (s *service) ListByUser(userID string) ([]Order, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    var result []Order
    for _, o := range s.store {
        if o.UserID == userID {
            result = append(result, o)
        }
    }
    return result, nil
}
