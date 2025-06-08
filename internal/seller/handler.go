package seller

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

// RegisterRoutes mounts seller endpoints under /sellers
func RegisterRoutes(r *gin.Engine, svc Service) {
    grp := r.Group("/sellers")

    // POST /sellers/register
    grp.POST("/register", func(c *gin.Context) {
        var s Seller
        if err := c.ShouldBindJSON(&s); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        if err := svc.Register(s); err != nil {
            c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusCreated, gin.H{"message": "seller registered"})
    })

    // GET /sellers/:id
    grp.GET("/:id", func(c *gin.Context) {
        id := c.Param("id")
        seller, err := svc.GetByID(id)
        if err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
            return
        }
        c.JSON(http.StatusOK, seller)
    })

    // GET /sellers (list all)
    grp.GET("", func(c *gin.Context) {
        all := svc.ListAll()
        c.JSON(http.StatusOK, all)
    })
}
