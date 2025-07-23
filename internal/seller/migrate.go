// in internal/seller/migrate.go
package seller

import "gorm.io/gorm"

func Migrate(db *gorm.DB) error {
    return db.AutoMigrate(&Seller{})
}
