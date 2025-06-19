package order

import (
	"errors"
	"sync"
	"time"

	"github.com/albus-droid/Capstone-Project-Backend/internal/event"
)

type inMemoryService struct {
	mu     sync.Mutex
	orders map[string]Order
}

func NewInMemoryService() Service {
	return &inMemoryService{
		orders: make(map[string]Order),
	}
}

func (s *inMemoryService) Create(o Order) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.orders[o.ID]; exists {
		return errors.New("order already exists")
	}

	o.CreatedAt = time.Now().Unix()
	s.orders[o.ID] = o

	// 🔥 Emit OrderPlaced
	go func() {
		event.Bus <- event.Event{
			Type: "OrderPlaced",
			Data: o,
		}
	}()

	return nil
}

func (s *inMemoryService) GetByID(id string) (*Order, error) {
	o, ok := s.orders[id]
	if !ok {
		return nil, errors.New("order not found")
	}
	return &o, nil
}

func (s *inMemoryService) ListByUser(userID string) ([]Order, error) {
	var list []Order
	for _, o := range s.orders {
		if o.UserID == userID {
			list = append(list, o)
		}
	}
	return list, nil
}

func (s *inMemoryService) Accept(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	order, ok := s.orders[id]
	if !ok {
		return errors.New("order not found")
	}

	// no status change, but still trigger event
	go func() {
		event.Bus <- event.Event{
			Type: "OrderAccepted",
			Data: order,
		}
	}()

	return nil
}
