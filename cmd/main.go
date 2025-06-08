package main

import (
  "github.com/gin-gonic/gin"
  "github.com/albus-droid/Capstone-Project-Backend/internal/user"
)

func main() {
  r := gin.Default()
  usvc := user.NewInMemoryService()
  user.RegisterRoutes(r, usvc)
  // later: register seller, listing, order here too
  r.Run(":8080")  // http://localhost:8080
}
