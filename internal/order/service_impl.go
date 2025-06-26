package order

import (
	"errors"
	"sync"
	"time"

	"github.com/albus-droid/Capstone-Project-Backend/internal/event"
)

// ─────────────────────────────────────────────────────────────────────────────
// In-memory implementation
// ─────────────────────────────────────────────────────────────────────────────

type inMemoryService struct {
	mu     sync.Mutex
	orders map[string]Order // keyed by Order.ID
}

func NewInMemoryService() Service {
	return &inMemoryService{
		orders: make(map[string]Order),
	}
}

// Create stores a new order and emits OrderPlaced.
func (s *inMemoryService) Create(o Order) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.orders[o.ID]; exists {
		return errors.New("order already exists")
	}

	o.CreatedAt = time.Now().Unix()
	s.orders[o.ID] = o

	go s.emit("OrderPlaced", o)
	return nil
}

// GetByID returns a copy of the order. Caller must check ownership.
func (s *inMemoryService) GetByID(id string) (*Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	o, ok := s.orders[id]
	if !ok {
		return nil, errors.New("order not found")
	}
	return &o, nil
}

// ListByUser returns all orders that belong to the given e-mail.
func (s *inMemoryService) ListByUser(userEmail string) ([]Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var list []Order
	for _, o := range s.orders {
		if o.UserEmail == userEmail {
			list = append(list, o)
		}
	}
	return list, nil
}

// Accept marks the order as accepted if caller owns it.
func (s *inMemoryService) Accept(id, callerEmail string) error {
	return s.updateStatus(id, callerEmail, "accepted", "OrderAccepted")
}

// Complete marks the order as completed if caller owns it.
func (s *inMemoryService) Complete(id, callerEmail string) error {
	return s.updateStatus(id, callerEmail, "completed", "OrderCompleted")
}

// ─────────────────────────────────────────────────────────────────────────────
// helpers
// ─────────────────────────────────────────────────────────────────────────────

func (s *inMemoryService) updateStatus(id, callerEmail, newStatus, eventType string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	o, ok := s.orders[id]
	if !ok {
		return errors.New("order not found")
	}
	if o.UserEmail != callerEmail {
		return errors.New("forbidden")
	}

	o.Status = newStatus
	s.orders[id] = o

	go s.emit(eventType, o)
	return nil
}

func (s *inMemoryService) emit(typ string, data Order) {
	event.Bus <- event.Event{
		Type: typ,
		Data: data,
	}
}
