package order

import (
    "gorm.io/gorm"
	"gorm.io/datatypes"
)

// Order is the GORM model for an order record
// swagger:model Order
// gorm.Model is not embedded so we control fields explicitly
type Order struct {
    ID         string         `json:"id" gorm:"type:uuid;primaryKey"`
    // UserID     string         `json:"userId" gorm:"type:uuid;not null;index"`
    UserEmail  string         `json:"user_email" gorm:"type:varchar(100);not null;index"`
    SellerID   string         `json:"sellerId" gorm:"type:uuid;not null;index"`
    ListingIDs datatypes.JSON `json:"listingIds" gorm:"type:jsonb;not null;default:'[]'"`
	Total      float64        `json:"total" gorm:"type:numeric;not null"`
    CreatedAt  int64          `json:"createdAt" gorm:"autoCreateTime"`
    Status     string         `json:"status" gorm:"type:varchar(20);not null;default:'pending'"`
    DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"` // optional soft-delete
}
