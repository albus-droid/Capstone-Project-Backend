package seller

import (
	"errors"
	"sort"
	"sync"
)

// Seller represents a marketplace seller.
type Seller struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// ─────────────────────────────────────────────────────────────────────────────
// In-memory implementation
// ─────────────────────────────────────────────────────────────────────────────

type inMemoryService struct {
	mu      sync.Mutex
	sellers map[string]Seller // keyed by Seller.ID
}

func NewInMemoryService() Service {
	return &inMemoryService{
		sellers: make(map[string]Seller),
	}
}

func (s *inMemoryService) Register(sl Seller) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.sellers[sl.ID]; exists {
		return errors.New("seller already exists")
	}
	s.sellers[sl.ID] = sl
	return nil
}

func (s *inMemoryService) GetByID(id string) (Seller, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	sl, ok := s.sellers[id]
	if !ok {
		return Seller{}, errors.New("seller not found")
	}
	return sl, nil
}

func (s *inMemoryService) ListAll() []Seller {
	s.mu.Lock()
	defer s.mu.Unlock()

	out := make([]Seller, 0, len(s.sellers))
	for _, v := range s.sellers {
		out = append(out, v)
	}
	// deterministic order helps tests
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}
