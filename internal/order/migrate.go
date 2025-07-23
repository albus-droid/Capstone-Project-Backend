// internal/order/migrate.go
package order

import "gorm.io/gorm"

func Migrate(db *gorm.DB) error {
    return db.AutoMigrate(&Order{})
}
