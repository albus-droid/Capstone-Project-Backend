package order

import (
    "errors"
    "time"
    "sort"

    "github.com/google/uuid"
    "gorm.io/gorm"
    "github.com/albus-droid/Capstone-Project-Backend/internal/event"
)

type postgresService struct {
    db *gorm.DB
}

func NewPostgresService(db *gorm.DB) Service {
    return &postgresService{db: db}
}

func (s *postgresService) Create(o *Order) error {
    // 1) assign a fresh ID & timestamp
    o.ID = uuid.NewString()
    o.CreatedAt = time.Now().Unix()
    o.Status = "pending"      // if not already set

    // 2) insert—UUID collisions are practically impossible, so no pre‑check needed
    if err := s.db.Create(o).Error; err != nil {
        return err
    }

    // 3) emit the OrderPlaced event
    go func(placed Order) {
        event.Bus <- event.Event{Type: "OrderPlaced", Data: placed}
    }(*o)

    return nil
}

func (s *postgresService) GetByID(id string) (*Order, error) {
    var o Order
    if err := s.db.First(&o, "id = ?", id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("order not found")
        }
        return nil, err
    }
    return &o, nil
}

func (s *postgresService) ListByUser(userEmail string) ([]Order, error) {
    var list []Order
    if err := s.db.Where("user_email = ?", userEmail).Find(&list).Error; err != nil {
        return nil, err
    }
    // keep same ordering as in-memory
    sort.Slice(list, func(i, j int) bool {
        return list[i].CreatedAt < list[j].CreatedAt
    })
    return list, nil
}

func (s *postgresService) Accept(id, callerEmail string) error {
    return s.updateStatus(id, callerEmail, "accepted", "OrderAccepted")
}

func (s *postgresService) Complete(id, callerEmail string) error {
    return s.updateStatus(id, callerEmail, "completed", "OrderCompleted")
}

func (s *postgresService) updateStatus(id, callerEmail, newStatus, eventType string) error {
    var o Order
    if err := s.db.First(&o, "id = ?", id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return errors.New("order not found")
        }
        return err
    }
    if o.UserEmail != callerEmail {
        return errors.New("forbidden")
    }

    if err := s.db.Model(&o).Update("status", newStatus).Error; err != nil {
        return err
    }
    o.Status = newStatus
    // emit event
    go func(ev event.Event) {
       event.Bus <- ev
    }(event.Event{Type: eventType, Data: o})
    return nil
}
