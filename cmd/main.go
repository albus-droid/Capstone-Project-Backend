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

	r.Run(":8080") // http://localhost:8080
}
