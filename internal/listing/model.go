package listing

import (
    "time"

    "gorm.io/gorm"
)

// Listing is the GORM model for an item listing
type Listing struct {
    ID          string         `json:"id" gorm:"type:uuid;primaryKey"`
    SellerID    string         `json:"sellerId" gorm:"type:uuid;not null;index"`
    Title       string         `json:"title" gorm:"type:varchar(200);not null"`
    Description string         `json:"description" gorm:"type:text"`
    Price       float64        `json:"price" gorm:"not null"`
    Available   bool           `json:"available" gorm:"default:true;not null"`
    CreatedAt   time.Time      `json:"createdAt" gorm:"autoCreateTime"`
    UpdatedAt   time.Time      `json:"updatedAt" gorm:"autoUpdateTime"`
    DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"` // soft‑delete
}
