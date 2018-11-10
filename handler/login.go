package handler

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/globalsign/mgo/bson"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"github.com/jaaaaason/hmblog/database"
	"github.com/jaaaaason/hmblog/structure"
)

// PostLogin handles the POST request for /admin/login
func PostLogin(c *gin.Context) {
	var login structure.Login
	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, errRes{
			Status:  http.StatusBadRequest,
			Message: "Bad request",
		})
		return
	}

	// check if user exists
	user, err := database.User(bson.M{
		"username": login.Username,
	})
	if err != nil {
		// user not exist
		if err == database.ErrNoUser {
			c.JSON(http.StatusNotFound, errRes{
				Status:  http.StatusNotFound,
				Message: "No user named " + login.Username,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, errRes{
			Status:  http.StatusInternalServerError,
			Message: "Internal server error",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(login.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, errRes{
			Status:  http.StatusBadRequest,
			Message: "Wrong password",
		})
		return
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := jwt.MapClaims{
		"exp": time.Now().Add(time.Second * tokenExp).Unix(),
		"iat": time.Now().Unix(),
	}
	claims["user_id"] = user.ID
	token.Claims = claims

	tokenString, err := token.SignedString([]byte(jwtSignKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errRes{
			Status:  http.StatusInternalServerError,
			Message: "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": tokenString,
		"token_type":   tokenType,
		"expires_in":   tokenExp,
	})
}
