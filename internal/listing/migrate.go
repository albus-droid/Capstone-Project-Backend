package listing

import "gorm.io/gorm"

// Migrate creates/updates the listings table to match the model
func Migrate(db *gorm.DB) error {
    return db.AutoMigrate(&Listing{})
}
