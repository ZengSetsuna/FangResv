package api

import (
	"FangResv/auth"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const secretKey = "dF453fEsEV3bjfnd29cFoLpq8432fn9O"

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// 解析 "Bearer <token>"
		authHeader = strings.TrimSpace(authHeader)
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			log.Println(tokenParts, len(tokenParts), tokenParts[0])
			c.Abort()
			return
		}
		tokenString := tokenParts[1]

		// 验证 Token
		status, payload, err := auth.ValidateToken(tokenString)
		if err != nil {
			if status == auth.TokenExpired {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			}
			c.Abort()
			return
		}
		userID := payload.UserID

		// 存入 Gin 上下文
		c.Set("user_id", int32(userID))
		c.Next()
	}
}
