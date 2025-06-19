package main

import (
	"fmt"

	"github.com/albus-droid/Capstone-Project-Backend/internal/event"
	"github.com/albus-droid/Capstone-Project-Backend/internal/listing"
	"github.com/albus-droid/Capstone-Project-Backend/internal/order"
	"github.com/albus-droid/Capstone-Project-Backend/internal/seller"
	"github.com/albus-droid/Capstone-Project-Backend/internal/user"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	// user routes
	usvc := user.NewInMemoryService()
	user.RegisterRoutes(r, usvc)

	// seller routes
	ssvc := seller.NewInMemoryService()
	seller.RegisterRoutes(r, ssvc)

	// Listing routes
	lsvc := listing.NewInMemoryService()
	listing.RegisterRoutes(r, lsvc)

	// Order
	osvc := order.NewInMemoryService()
	order.RegisterRoutes(r, osvc)

	startNotificationListener()
	r.Run(":8080") // http://localhost:8080
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
				fmt.Printf("📬 Notify user %s that order %s was accepted\n", order.UserID, order.ID)
			}
		}
	}()
}
