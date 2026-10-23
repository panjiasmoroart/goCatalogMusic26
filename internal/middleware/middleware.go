package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/panjiasmoroart/goCatalogMusic26/internal/configs"
	"github.com/panjiasmoroart/goCatalogMusic26/pkg/jwt"
)

func AuthMiddleware() gin.HandlerFunc {
	secretKey := configs.Get().Service.SecretKey
	return func(c *gin.Context) {
		header := c.Request.Header.Get("Authorization")

		header = strings.TrimSpace(header)
		if header == "" {
			// c.AbortWithError(http.StatusUnauthorized, errors.New("missing token"))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "missing token",
			})
			return
		}

		userID, username, err := jwt.ValidateToken(header, secretKey)
		if err != nil {
			// c.AbortWithError(http.StatusUnauthorized, err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			return
		}
		c.Set("userID", userID)
		c.Set("username", username)
		c.Next()
	}
}
