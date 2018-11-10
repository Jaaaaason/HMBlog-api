package handler

import (
	"errors"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const (
	jwtSignKey = "secret"
	tokenType  = "bearer"
	tokenExp   = 86400 // 1 Day, 86400 second
)

// JWTMiddleware the middleware for verifying jwt token
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, errRes{
				Status:  http.StatusUnauthorized,
				Message: "JWT token required",
			})

			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString,
			func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("Can't parse JWT token")
				}

				return []byte(jwtSignKey), nil
			})
		if err != nil {
			v, _ := err.(*jwt.ValidationError)
			if v.Errors == jwt.ValidationErrorExpired {
				// jwt token is expired
				c.JSON(http.StatusUnauthorized, errRes{
					Status:  http.StatusUnauthorized,
					Message: "JWT token is expired",
				})
			} else {
				c.JSON(http.StatusInternalServerError, errRes{
					Status:  http.StatusInternalServerError,
					Message: "Internal server error",
				})
			}

			c.Abort()
			return
		}

		// invalid token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, errRes{
				Status:  http.StatusUnauthorized,
				Message: "Invalid JWT token",
			})

			c.Abort()
			return
		}

		c.Set("user_id", claims["user_id"])
		c.Next()
	}
}
