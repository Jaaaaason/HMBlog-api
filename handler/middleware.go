package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const (
	jwtSignKey = "secret"
	tokenType  = "bearer"
	tokenExp   = 86400 // 1 Day, 86400 seconds
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

		tokenStrs := strings.Split(tokenString, " ")
		if len(tokenStrs) != 2 || tokenStrs[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, errRes{
				Status:  http.StatusUnauthorized,
				Message: "Invalid jwt token",
			})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenStrs[1],
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
				c.JSON(http.StatusUnauthorized, errRes{
					Status:  http.StatusUnauthorized,
					Message: "Invalid jwt token",
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

// CORSMiddleware the middleware for cors
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTION")

		c.Next()
	}
}
