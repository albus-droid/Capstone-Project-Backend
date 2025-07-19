package listing

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

// Service defines the methods our handlers expect.
// Make sure your implementation has Create, GetByID, ListBySeller, ListAll,
// Update(id string, l Listing) error, and Delete(id string) error.
type Service interface {
    Create(l Listing) error
    GetByID(id string) (Listing, error)
    ListBySeller(sellerID string) ([]Listing, error)
    ListAll() []Listing
    Update(id string, l Listing) error
    Delete(id string) error
}

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

    // PUT /listings/:id  — update an existing listing
    grp.PUT("/:id", func(c *gin.Context) {
        id := c.Param("id")
        var l Listing
        if err := c.ShouldBindJSON(&l); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        if err := svc.Update(id, l); err != nil {
            // you can customize error handling based on your svc.Update error
            c.JSON(http.StatusNotFound, gin.H{"error": "not found or unable to update"})
            return
        }
        c.JSON(http.StatusOK, gin.H{"message": "listing updated"})
    })

    // DELETE /listings/:id  — remove a listing
    grp.DELETE("/:id", func(c *gin.Context) {
        id := c.Param("id")
        if err := svc.Delete(id); err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "not found or unable to delete"})
            return
        }
        c.Status(http.StatusNoContent)
    })
}
