package seller

import (
    "net/http"
    "time"

    "github.com/albus-droid/Capstone-Project-Backend/internal/auth"
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v4"
)

// RegisterRoutes mounts seller endpoints under /sellers
func RegisterRoutes(r *gin.Engine, svc Service, store auth.Store) {
    grp := r.Group("/sellers")

    // Register
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

    // Login → returns JWT
    grp.POST("/login", func(c *gin.Context) {
        var creds struct {
            Email    string `json:"email"`
            Password string `json:"password"`
        }
        if err := c.ShouldBindJSON(&creds); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        // Authenticate against your seller service
        seller, err := svc.Authenticate(creds.Email, creds.Password)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
            return
        }

        // Create JWT with seller’s email (or ID) as subject
        token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
            "sub": seller.Email,
            "exp": time.Now().Add(24 * time.Hour).Unix(),
        })

        signed, err := token.SignedString(auth.Secret())
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "could not sign token"})
            return
        }

        // compute TTL from "exp" claim
        exp := time.Unix(token.Claims.(jwt.MapClaims)["exp"].(int64), 0)
        ttl := time.Until(exp)
        // persist it in Redis
        if err := store.Save(c.Request.Context(), signed, ttl); err != nil {
            c.Error(err) // log but don’t block
        }

        c.JSON(http.StatusOK, gin.H{"token": signed})
    })

    // Fetch a seller by ID
    grp.GET("/:id", func(c *gin.Context) {
        id := c.Param("id")
        seller, err := svc.GetByID(id)
        if err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
            return
        }
        c.JSON(http.StatusOK, seller)
    })

    // List all sellers
    grp.GET("", func(c *gin.Context) {
        all := svc.ListAll()
        c.JSON(http.StatusOK, all)
    })
}
