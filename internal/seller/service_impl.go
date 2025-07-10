package seller

import (
	"errors"
	"sort"
	"sync"
	"golang.org/x/crypto/bcrypt"
)

// ─────────────────────────────────────────────────────────────────────────────
// In-memory implementation
// ─────────────────────────────────────────────────────────────────────────────

type inMemoryService struct {
	mu      sync.Mutex
	sellers map[string]Seller // keyed by Seller.Email
}

func NewInMemoryService() Service {
	return &inMemoryService{
		sellers: make(map[string]Seller),
	}
}

func (s *inMemoryService) Register(sl Seller) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.sellers[sl.Email]; exists {
		return errors.New("seller already exists")
	}

	// Assign new UUID
    sl.ID := uuid.New().String()
    
	// hash password
	h, err := bcrypt.GenerateFromPassword([]byte(sl.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	sl.Password = string(h)
	sl.Verified = false
	s.sellers[sl.Email] = sl
	return nil
}

func (s *inMemoryService) GetByID(id string) (*Seller, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, sl := range s.sellers {
        if sl.ID == id {
            return &sl, nil
        }
    }
	return nil, errors.New("seller not found")
}

func (s *inMemoryService) ListAll() []Seller {
	s.mu.Lock()
	defer s.mu.Unlock()

	out := make([]Seller, 0, len(s.sellers))
	for _, v := range s.sellers {
		out = append(out, v)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

func (s *inMemoryService) Authenticate(email, pw string) (*Seller, error) {
  u, err := s.GetByEmail(email)
  if err != nil {
    return nil, err
  }
  if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pw)); err != nil {
    return nil, errors.New("invalid credentials")
  }
  return u, nil
}

func (s *inMemoryService) GetByEmail(email string) (*Seller, error) {
 s.mu.Lock()
    defer s.mu.Unlock()

    sl, ok := s.sellers[email]
    if !ok {
        return nil, errors.New("seller not found")
    }
    return &sl, nil
}
