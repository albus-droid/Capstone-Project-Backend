// internal/auth/middleware.go
package auth   // new package so every module can import it

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET")) // never hard-code

const ctxEmailKey = "userEmail"

// Middleware validates the JWT and stores the user’s e-mail in Gin’s context.
func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
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
		c.Set(ctxEmailKey, claims["sub"].(string)) // store e-mail
		c.Next()
	}
}
