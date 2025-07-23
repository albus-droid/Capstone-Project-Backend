// internal/seller/postgres_service.go
package seller

import (
    "errors"
    "sort"

    "github.com/google/uuid"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
)

type postgresService struct {
    db *gorm.DB
}

func NewPostgresService(db *gorm.DB) Service {
    return &postgresService{db: db}
}

func (s *postgresService) Register(sl Seller) error {
    // check for existing email
    var cnt int64
    if err := s.db.Model(&Seller{}).
                  Where("email = ?", sl.Email).
                  Count(&cnt).Error; err != nil {
        return err
    }
    if cnt > 0 {
        return errors.New("seller already exists")
    }

    // assign UUID & hash pw
    sl.ID = uuid.New().String()
    h, err := bcrypt.GenerateFromPassword([]byte(sl.Password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    sl.Password = string(h)
    sl.Verified = false

    return s.db.Create(&sl).Error
}

func (s *postgresService) GetByID(id string) (*Seller, error) {
    var sl Seller
    if err := s.db.First(&sl, "id = ?", id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("seller not found")
        }
        return nil, err
    }
    return &sl, nil
}

func (s *postgresService) ListAll() []Seller {
    var all []Seller
    s.db.Find(&all)
    // keep your in‑memory sort behavior
    sort.Slice(all, func(i, j int) bool {
        return all[i].ID < all[j].ID
    })
    return all
}

func (s *postgresService) Authenticate(email, pw string) (*Seller, error) {
    sl, err := s.GetByEmail(email)
    if err != nil {
        return nil, err
    }
    if bcrypt.CompareHashAndPassword([]byte(sl.Password), []byte(pw)) != nil {
        return nil, errors.New("invalid credentials")
    }
    return sl, nil
}

func (s *postgresService) GetByEmail(email string) (*Seller, error) {
    var sl Seller
    if err := s.db.First(&sl, "email = ?", email).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("seller not found")
        }
        return nil, err
    }
    return &sl, nil
}
