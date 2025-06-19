package order

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

// RegisterRoutes mounts order endpoints under /orders
func RegisterRoutes(r *gin.Engine, svc Service) {
    grp := r.Group("/orders")

    // POST /orders
    grp.POST("", func(c *gin.Context) {
        var o Order
        if err := c.ShouldBindJSON(&o); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        if err := svc.Create(o); err != nil {
            c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusCreated, gin.H{"message": "order created"})
    })

    // GET /orders/:id
    grp.GET("/:id", func(c *gin.Context) {
        id := c.Param("id")
        o, err := svc.GetByID(id)
        if err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
            return
        }
        c.JSON(http.StatusOK, o)
    })

    // GET /orders?userId=...
    grp.GET("", func(c *gin.Context) {
        userID := c.Query("userId")
        if userID == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "userId query param required"})
            return
        }
        list, _ := svc.ListByUser(userID)
        c.JSON(http.StatusOK, list)
    })

    grp.PATCH("/:id/accept", func(c *gin.Context) {
    id := c.Param("id")
    if err := svc.Accept(id); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "order accepted"})
})

}
