package middleware

import (
	"net/http"
	"server-go/internal/config"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusOK, gin.H{"code": 401, "message": "Authorization header required"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusOK, gin.H{"code": 401, "message": "Invalid authorization header format"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusOK, gin.H{"code": 401, "message": "Invalid token"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			idFloat, ok := claims["id"].(float64)
			if !ok {
				c.JSON(http.StatusOK, gin.H{"code": 401, "message": "Invalid token claims"})
				c.Abort()
				return
			}
			userID := int(idFloat)
			username, _ := claims["username"].(string)
			role, _ := claims["role"].(string)
			realName := ""
			if rn, ok := claims["realName"].(string); ok {
				realName = rn
			}

			c.Set("userID", userID)
			c.Set("username", username)
			c.Set("role", role)
			c.Set("realName", realName)
		}

		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			c.JSON(http.StatusOK, gin.H{"code": 403, "message": "Admin access required"})
			c.Abort()
			return
		}
		c.Next()
	}
}
