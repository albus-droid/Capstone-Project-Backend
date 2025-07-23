package listing

import (
    "errors"
    "sort"

    "github.com/google/uuid"
    "gorm.io/gorm"
)

// postgresService persists listings in Postgres via GORM.
type postgresService struct {
    db *gorm.DB
}

// NewPostgresService returns a Listing Service backed by Postgres.
func NewPostgresService(db *gorm.DB) Service {
    return &postgresService{db}
}

func (s *postgresService) Create(l *Listing) error {
    l.ID = uuid.NewString()
    return s.db.Create(l).Error
}

func (s *postgresService) GetByID(id string) (*Listing, error) {
    var l Listing
    if err := s.db.First(&l, "id = ?", id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("listing not found")
        }
        return nil, err
    }
    return &l, nil
}

func (s *postgresService) ListBySeller(sellerID string) ([]Listing, error) {
    var out []Listing
    if err := s.db.Where("seller_id = ?", sellerID).Find(&out).Error; err != nil {
        return nil, err
    }
    sort.Slice(out, func(i, j int) bool {
        return out[i].CreatedAt.Before(out[j].CreatedAt)
    })
    return out, nil
}

// If you need a ListAll (not in your in‑memory), you can add:
func (s *postgresService) ListAll() []Listing {
    var all []Listing
    s.db.Find(&all)
    sort.Slice(all, func(i, j int) bool {
        return all[i].CreatedAt.Before(all[j].CreatedAt)
    })
    return all
}

func (s *postgresService) Update(id string, l Listing) error {
    l.ID = id
    return s.db.Save(&l).Error
}

func (s *postgresService) Delete(id string) error {
    return s.db.Delete(&Listing{}, "id = ?", id).Error
}
