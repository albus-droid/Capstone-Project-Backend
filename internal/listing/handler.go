package listing

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

// RegisterRoutes mounts listing endpoints under /listings
func RegisterRoutes(r *gin.Engine, svc Service) {
    grp := r.Group("/listings")

    // POST /listings
    grp.POST("", func(c *gin.Context) {
        var l Listing
        if err := c.ShouldBindJSON(&l); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        if err := svc.Create(l); err != nil {
            c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusCreated, gin.H{"message": "listing created"})
    })

    // GET /listings/:id
    grp.GET("/:id", func(c *gin.Context) {
        id := c.Param("id")
        l, err := svc.GetByID(id)
        if err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
            return
        }
        c.JSON(http.StatusOK, l)
    })

    // GET /listings?sellerId=...
    grp.GET("", func(c *gin.Context) {
        if sellerID := c.Query("sellerId"); sellerID != "" {
            list, _ := svc.ListBySeller(sellerID)
            c.JSON(http.StatusOK, list)
            return
        }
        // fallback: list all
        c.JSON(http.StatusOK, svc.ListAll())
    })
}
