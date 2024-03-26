package middlewares

import (
	"fmt"
	"net/http"
	"new_projects/initializers"
	"new_projects/models"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authoriztion")

		if !(strings.HasPrefix(tokenString, "Bearer ")) {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "the token has no bearer prefix"})
			return
		}

		tokenStr := strings.TrimPrefix(tokenString, "Bearer ")

		// decode or validate it
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(os.Getenv("SECRET")), nil
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// Check exp
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				c.JSON(http.StatusGatewayTimeout, gin.H{"Error": "token Expired"})
				return
			}
			// fined the user with token
			var user models.User
			initializers.DB.First(&user, claims["sub"])

			if user.ID == 0 {
				c.JSON(http.StatusBadRequest, gin.H{"Error": "Failed to get the user with the token"})
				return
			}

			// attach the request
			c.Set("user_id", user.ID)
			c.Set("role", user.Role)
			fmt.Println("user data", user.Role)

			// continue
			c.Next()
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "user unauthenticated"})
			return
		}
	}
}

func Authoriztion() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleInterface, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Role not found"})
			return
		}

		roleString, ok := roleInterface.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Invalid error"})
			return
		}

		method := c.Request.Method
		if models.Role(roleString) == models.UserRole {
			if method != http.MethodGet {
				c.JSON(http.StatusUnauthorized, gin.H{"data": "unauthorized"})
				return
			}
		}
		c.Next()
	}
}
