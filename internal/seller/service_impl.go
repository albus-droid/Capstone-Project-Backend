package seller

import "errors"

type inMemoryService struct {
    sellers map[string]Seller
}

// NewInMemoryService returns a new Seller service
func NewInMemoryService() Service {
    return &inMemoryService{
        sellers: make(map[string]Seller),
    }
}

func (s *inMemoryService) Register(seller Seller) error {
    if _, exists := s.sellers[seller.ID]; exists {
        return errors.New("seller already exists")
    }
    s.sellers[seller.ID] = seller
    return nil
}

func (s *inMemoryService) GetByID(id string) (*Seller, error) {
    seller, ok := s.sellers[id]
    if !ok {
        return nil, errors.New("seller not found")
    }
    return &seller, nil
}

func (s *inMemoryService) ListAll() []Seller {
    list := make([]Seller, 0, len(s.sellers))
    for _, v := range s.sellers {
        list = append(list, v)
    }
    return list
}
