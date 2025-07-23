package user

import (
	"net/http"
	"time"

	"github.com/albus-droid/Capstone-Project-Backend/internal/auth"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// RegisterRoutes mounts user endpoints under /users
func RegisterRoutes(r *gin.Engine, svc Service) {
	g := r.Group("/users")

	// Register
	g.POST("/register", func(c *gin.Context) {
		var u User
		if err := c.ShouldBindJSON(&u); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := svc.Register(u); err != nil {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "registered"})
	})

	// Login → returns JWT
	g.POST("/login", func(c *gin.Context) {
		var creds struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&creds); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		u, err := svc.Authenticate(creds.Email, creds.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": u.Email,
			"exp": time.Now().Add(24 * time.Hour).Unix(),
		})
		ts, err := token.SignedString(auth.Secret()) // 🔒 use shared secret
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not sign token"})
			return
		}

        // save it in Redis
        if err := ts.Save(c.Request.Context(), tsStr, time.Until(exp)); err != nil {
            // log but don’t block response
            c.Error(err)
        }

		c.JSON(http.StatusOK, gin.H{"token": ts})
	})

	// Profile (protected)
	g.GET("/profile", auth.Middleware(), func(c *gin.Context) {
		email := c.GetString(string(auth.CtxEmailKey))
		u, err := svc.GetByEmail(email)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusOK, u)
	})
}
