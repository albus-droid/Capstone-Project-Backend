package listing

import {
    "errors"
	"github.com/google/uuid"
}

type inMemoryService struct {
    items map[string]Listing
}

func NewInMemoryService() Service {
    return &inMemoryService{
        items: make(map[string]Listing),
    }
}

func (s *inMemoryService) Create(l Listing) error {

    l.ID = uuid.New().String()

    if _, exists := s.items[l.ID]; exists {
        return errors.New("listing already exists")
    }
    s.items[l.ID] = l
    return nil
}

func (s *inMemoryService) GetByID(id string) (*Listing, error) {
    l, ok := s.items[id]
    if !ok {
        return nil, errors.New("listing not found")
    }
    return &l, nil
}

func (s *inMemoryService) ListBySeller(sellerID string) ([]Listing, error) {
    var result []Listing
    for _, l := range s.items {
        if l.SellerID == sellerID {
            result = append(result, l)
        }
    }
    return result, nil
}

func (s *inMemoryService) ListAll() []Listing {
    var all []Listing
    for _, l := range s.items {
        all = append(all, l)
    }
    return all
}
