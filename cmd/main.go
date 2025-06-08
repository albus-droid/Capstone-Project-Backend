package main

import (
  "github.com/gin-gonic/gin"
  "github.com/albus-droid/Capstone-Project-Backend/internal/user"
  "github.com/albus-droid/Capstone-Project-Backend/internal/seller"

)

func main() {
  r := gin.Default()
 // user routes
    usvc := user.NewInMemoryService()
    user.RegisterRoutes(r, usvc)

    // seller routes
    ssvc := seller.NewInMemoryService()
    seller.RegisterRoutes(r, ssvc)
  r.Run(":8080")  // http://localhost:8080
}
