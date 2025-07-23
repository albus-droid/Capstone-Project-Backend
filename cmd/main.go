package main

import (
	"fmt"
	"github.com/albus-droid/Capstone-Project-Backend/internal/event"
	"github.com/albus-droid/Capstone-Project-Backend/internal/listing"
	"github.com/albus-droid/Capstone-Project-Backend/internal/order"
	"github.com/albus-droid/Capstone-Project-Backend/internal/seller"
	"github.com/albus-droid/Capstone-Project-Backend/internal/user"
	"github.com/albus-droid/Capstone-Project-Backend/internal/db"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	db := db.Init()

	// user routes
	user.Migrate(db) // optional for dev
	usvc := user.NewPostgresService(db)
	user.RegisterRoutes(r, usvc)

	// seller routes
	seller.Migrate(db) // optional for dev
	ssvc := seller.NewPostgresService(db)
	seller.RegisterRoutes(r, ssvc)

	// Listing routes
	listing.Migrate(db) // optional for dev
	lsvc := listing.NewPostgresService(db)
	listing.RegisterRoutes(r, lsvc)

	// Order
	order.Migrate(db) // optional for dev
	osvc := order.NewPostgresService(db)
	order.RegisterRoutes(r, osvc)

	startNotificationListener()
	r.Run(":8000") // http://localhost:8080
}

func startNotificationListener() {
	go func() {
		for e := range event.Bus {
			switch e.Type {

			case "OrderPlaced":
				order := e.Data.(order.Order)
				fmt.Printf("📦 Notify seller %s of new order %s\n", order.SellerID, order.ID)

			case "OrderAccepted":
				order := e.Data.(order.Order)
				fmt.Printf("📬 Notify user %s that order %s was accepted\n", order.UserEmail, order.ID)
			}
		}
	}()
}
