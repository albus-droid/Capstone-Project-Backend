package order

import (
	"net/http"
	"log"
	"encoding/json"
	"errors"
	"github.com/albus-droid/Capstone-Project-Backend/internal/auth"
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

var ErrOrderAlreadyExists = errors.New("order already exists")

func RegisterRoutes(r *gin.Engine, svc Service) {
	grp := r.Group("/orders")
	grp.Use(auth.Middleware())

	// ─────────────────────────────────────────────────────────────
	// POST /orders – create an order
	// ─────────────────────────────────────────────────────────────
	grp.POST("", func(c *gin.Context) {
		var payload struct {
			ListingIDs []string `json:"listingIds"`
			SellerID   string   `json:"sellerId"`
			Total      float64  `json:"total"`
		}
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		raw, _ := json.Marshal(payload.ListingIDs)
		email := c.GetString(string(auth.CtxEmailKey))
		o := &Order{
			UserEmail:  email,
			SellerID:   payload.SellerID,
			ListingIDs: datatypes.JSON(raw),
			Total:      payload.Total,
		}
		if err := svc.Create(o); err != nil {
    		switch {
    		case errors.Is(err, ErrOrderAlreadyExists):
        		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
    		default:
        		log.Printf("create order failed: %v", err)
        		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    		}
    		return
    	}
    	c.JSON(http.StatusCreated, o) // return the order (or a message if you prefer)
	})

	// ─────────────────────────────────────────────────────────────
	// GET /orders/:id – only owner can fetch
	// ─────────────────────────────────────────────────────────────
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

	// ─────────────────────────────────────────────────────────────
	// GET /orders – list my orders
	// ─────────────────────────────────────────────────────────────
	grp.GET("", func(c *gin.Context) {
		email := c.GetString(string(auth.CtxEmailKey))
		list, _ := svc.ListByUser(email)
		c.JSON(http.StatusOK, list)
	})

	// ─────────────────────────────────────────────────────────────
	// PATCH /orders/:id/accept
	// ─────────────────────────────────────────────────────────────
	grp.PATCH("/:id/accept", func(c *gin.Context) {
		id := c.Param("id")
		email := c.GetString(string(auth.CtxEmailKey))
		if err := svc.Accept(id, email); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "order accepted"})
	})

	// ─────────────────────────────────────────────────────────────
	// PATCH /orders/:id/complete
	// ─────────────────────────────────────────────────────────────
	grp.PATCH("/:id/complete", func(c *gin.Context) {
		id := c.Param("id")
		email := c.GetString(string(auth.CtxEmailKey))
		if err := svc.Complete(id, email); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "order completed"})
	})
}
