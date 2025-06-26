// internal/order/routes.go
package order

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/albus-droid/Capstone-Project-Backend/internal/auth"
)

// RegisterRoutes mounts order endpoints under /orders
func RegisterRoutes(r *gin.Engine, svc Service) {
	grp := r.Group("/orders")
	grp.Use(auth.Middleware())                  // 🔒 every route below is authenticated

	// POST /orders  – create order FOR THE CALLER
	grp.POST("", func(c *gin.Context) {
		var payload struct {
			Items []string `json:"items"` // whatever you need – NO userId here
		}
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		email := c.GetString(auth.CtxEmailKey)   // read from context
		o := Order{
			UserEmail: email,
			Items:     payload.Items,
		}

		if err := svc.Create(o); err != nil {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "order created"})
	})

	// GET /orders/:id – only owner can fetch
	grp.GET("/:id", func(c *gin.Context) {
		id := c.Param("id")
		o, err := svc.GetByID(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		if o.UserEmail != c.GetString(auth.CtxEmailKey) {
			c.JSON(http.StatusForbidden, gin.H{"error": "no access"})
			return
		}
		c.JSON(http.StatusOK, o)
	})

	// GET /orders – list *my* orders, ignore query params
	grp.GET("", func(c *gin.Context) {
		email := c.GetString(auth.CtxEmailKey)
		list, _ := svc.ListByUser(email)
		c.JSON(http.StatusOK, list)
	})

	// PATCH /orders/:id/accept – seller or admin could call; still owner-check
	grp.PATCH("/:id/accept", func(c *gin.Context) {
		id := c.Param("id")

		// Example: svc.Accept ensures the caller owns / is allowed
		if err := svc.Accept(id, c.GetString(auth.CtxEmailKey)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "order accepted"})
	})
}
