package order

import (
	"net/http"
	"time"

	"github.com/albus-droid/Capstone-Project-Backend/internal/auth"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes mounts order endpoints under /orders.
func RegisterRoutes(r *gin.Engine, svc Service) {
	grp := r.Group("/orders")
	grp.Use(auth.Middleware()) // all routes below are authenticated

	// POST /orders – create an order for the caller
	grp.POST("", func(c *gin.Context) {
		var payload struct {
			ListingIDs []string `json:"listingIds"`        // items in the cart
			SellerID   string   `json:"sellerId"`          // seller fulfilling it
			Total      float64  `json:"total"`             // client-calculated total
		}
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		email := c.GetString(string(auth.CtxEmailKey))
		order := Order{
			UserEmail:  email,               // owner
			SellerID:   payload.SellerID,
			ListingIDs: payload.ListingIDs,
			Total:      payload.Total,
			Status:     "pending",
			CreatedAt:  time.Now().Unix(),
		}

		if err := svc.Create(order); err != nil {
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
		if o.UserEmail != c.GetString(string(auth.CtxEmailKey)) {
			c.JSON(http.StatusForbidden, gin.H{"error": "no access"})
			return
		}
		c.JSON(http.StatusOK, o)
	})

	// GET /orders – list *my* orders
	grp.GET("", func(c *gin.Context) {
		email := c.GetString(string(auth.CtxEmailKey))
		list, _ := svc.ListByUser(email)
		c.JSON(http.StatusOK, list)
	})

	// PATCH /orders/:id/accept – owner (or privileged role) accepts
	grp.PATCH("/:id/accept", func(c *gin.Context) {
		id := c.Param("id")
		email := c.GetString(string(auth.CtxEmailKey))

		if err := svc.Accept(id, email); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "order accepted"})
	})
}
// PATCH /orders/:id/complete – owner (or privileged role) completes
func (svc Service) Complete(c *gin.Context) {
	id := c.Param("id")
	email := c.GetString(string(auth.CtxEmailKey))

	if err := svc.Complete(id, email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "order completed"})
}