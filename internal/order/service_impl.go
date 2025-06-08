package order

import (
    "errors"
    "time"
)

type inMemoryService struct {
    orders map[string]Order
}

func NewInMemoryService() Service {
    return &inMemoryService{
        orders: make(map[string]Order),
    }
}

func (s *inMemoryService) Create(o Order) error {
    if _, exists := s.orders[o.ID]; exists {
        return errors.New("order already exists")
    }
    o.CreatedAt = time.Now().Unix()
    s.orders[o.ID] = o
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
