package user

import (
  "net/http"
  "github.com/gin-gonic/gin"
  "github.com/golang-jwt/jwt/v4"
  "time"
)

var jwtSecret = []byte("replace-with-secure-secret")

// RegisterRoutes mounts user endpoints under /users
func RegisterRoutes(r *gin.Engine, svc Service) {
  g := r.Group("/users")

  // 1) Register
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

  // 2) Login → returns JWT
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
    // build token
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
      "sub": u.Email,
      "exp": time.Now().Add(24 * time.Hour).Unix(),
    })
    ts, err := token.SignedString(jwtSecret)
    if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": "could not sign token"})
      return
    }
    c.JSON(http.StatusOK, gin.H{"token": ts})
  })

  // 3) Profile (protected)
  g.GET("/profile", authMiddleware(), func(c *gin.Context) {
    email := c.GetString("userEmail")
    u, err := svc.GetByEmail(email)
    if err != nil {
      c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
      return
    }
    c.JSON(http.StatusOK, u)
  })
}

// simple JWT auth middleware
func authMiddleware() gin.HandlerFunc {
  return func(c *gin.Context) {
    auth := c.GetHeader("Authorization")
    if len(auth) < 7 || auth[:7] != "Bearer " {
      c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
      return
    }
    tok, err := jwt.Parse(auth[7:], func(t *jwt.Token) (interface{}, error) {
      if t.Method != jwt.SigningMethodHS256 {
        return nil, jwt.ErrSignatureInvalid
      }
      return jwtSecret, nil
    })
    if err != nil || !tok.Valid {
      c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
      return
    }
    claims := tok.Claims.(jwt.MapClaims)
    c.Set("userEmail", claims["sub"].(string))
    c.Next()
  }
}
