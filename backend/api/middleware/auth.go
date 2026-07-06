package middleware

import (
	"net/http"
	"strings"

	"blockchain/backend/service"
	"blockchain/backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		if !strings.HasPrefix(tokenStr, "Bearer ") {
			c.JSON(http.StatusOK, utils.Response{Code: 1002, Message: "未登录"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(strings.TrimPrefix(tokenStr, "Bearer "), func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(service.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusOK, utils.Response{Code: 1002, Message: "Token 无效"})
			c.Abort()
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Set("userID", claims["userID"])
		c.Set("username", claims["username"])
		c.Next()
	}
}
