package seller

import (
    "time"

    "gorm.io/gorm"
)

type Seller struct {
    ID        string         `json:"id" gorm:"type:uuid;primaryKey"`
    Name      string         `json:"name" gorm:"type:varchar(100);not null"`
    Email     string         `json:"email" gorm:"type:varchar(100);uniqueIndex;not null"`
    Password  string         `json:"password" gorm:"not null"`       // store a bcrypt hash
    Phone     string         `json:"phone" gorm:"type:varchar(20);not null"`
    Verified  bool           `json:"verified" gorm:"default:false"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`                // optional soft‑delete
}

