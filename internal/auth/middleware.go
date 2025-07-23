package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type ctxKey string

const CtxEmailKey ctxKey = "userEmail"

func Middleware(ts *auth.RedisStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

        raw := auth[7:]
        tok, err := jwt.Parse(raw, func(t *jwt.Token) (interface{}, error) {
            if t.Method != jwt.SigningMethodHS256 {
                return nil, jwt.ErrSignatureInvalid
            }
            return Secret(), nil
        })
        if err != nil || !tok.Valid {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
            return
        }

        // **New**: check Redis that the token hasn’t been revoked/expired
        exists, err := ts.Exists(c.Request.Context(), raw)
        if err != nil {
            // Redis error—fail closed
            c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "session store error"})
            return
        }
        if !exists {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token revoked or expired"})
            return
        }

        claims := tok.Claims.(jwt.MapClaims)
        c.Set(string(CtxEmailKey), claims["sub"].(string))
        c.Next()
    }
}
