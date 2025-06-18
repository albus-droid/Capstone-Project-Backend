package main

import (
	"github.com/albus-droid/Capstone-Project-Backend/internal/listing"
	"github.com/albus-droid/Capstone-Project-Backend/internal/order"
	"github.com/albus-droid/Capstone-Project-Backend/internal/seller"
	"github.com/albus-droid/Capstone-Project-Backend/internal/user"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	bus := events.NewBus()

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

    // subscribe to OrderPlacedEvent
    bus.Subscribe("OrderPlacedEvent", func(e events.Event) {
        evt := e.(events.OrderPlacedEvent)
        seller, found := sellerSvc.GetByID(evt.SellerID)
        if !found {
            log.Printf("❌ Seller %s not found for order %s", evt.SellerID, evt.OrderID)
            return
        }
        log.Printf("📦 Notify seller %s (%s): new order %s from %s", seller.Name, seller.Email, evt.OrderID, evt.UserEmail)
    })

	r.Run(":8080") // http://localhost:8080
}
