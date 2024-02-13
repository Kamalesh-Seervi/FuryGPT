package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kamalesh-seervi/simpleGPT/service"
)

var jwtKey = []byte(utils.GetSecretKey()) // Replace with your actual secret key

// Claims struct to store user details in the JWT
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			// Retrieve user details from the database
			user, err := service.GetUserByUsername(claims.Username)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
				c.Abort()
				return
			}

			// Set user details in the context
			c.Set("user", user)
			c.Next()
			return
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
	}
}