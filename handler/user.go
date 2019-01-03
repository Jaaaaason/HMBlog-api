package handler

import (
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"github.com/jaaaaason/hmblog/database"
	"github.com/jaaaaason/hmblog/structure"
	validator "gopkg.in/go-playground/validator.v8"
)

// UpdateUser handles the PUT and PATCH request
// for url path "/admin/user/:id"
func UpdateUser(c *gin.Context) {
	// parse object id from url path
	if !bson.IsObjectIdHex(c.Param("id")) {
		c.JSON(http.StatusBadRequest, errRes{
			Status:  http.StatusBadRequest,
			Message: "Invaild id",
		})
		return
	}

	oid := bson.ObjectIdHex(c.Param("id"))

	user, err := database.User(bson.M{
		"_id": oid,
	})
	if err != nil {
		if err == database.ErrNoUser {
			c.JSON(http.StatusNotFound, errRes{
				Status:  http.StatusNotFound,
				Message: "No such user",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, errRes{
			Status:  http.StatusInternalServerError,
			Message: "Internal server error",
		})
		return
	}

	if c.Request.Method == "PUT" {
		// for PUT request, use a new user struct,
		// binding with the request body, so the category
		// will be exactly the same as request body,
		// the value of some fields that doesn't provide will be
		// the zero value
		user = structure.User{}
		err = c.ShouldBindJSON(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, errRes{
				Status:  http.StatusBadRequest,
				Message: "Bad request",
			})
			return
		}
	} else if c.Request.Method == "PATCH" {
		// for PATCH request, use the origin user just got before,
		// binding with request body, so the value of some fields that
		// doesn't provide will not change
		err = c.ShouldBindJSON(&user)
		if err != nil {
			_, ok := err.(validator.ValidationErrors)
			if !ok {
				c.JSON(http.StatusBadRequest, errRes{
					Status:  http.StatusBadRequest,
					Message: "Bad request",
				})
				return
			}
		}
	}

	// set id zero value to omit it
	user.ID = nil

	// trim space
	user.Username = strings.TrimSpace(user.Username)
	if user.Username == "" {
		// empty category name
		c.JSON(http.StatusBadRequest, errRes{
			Status:  http.StatusBadRequest,
			Message: "Username shouldn't be just some whitespace",
		})
		return
	}

	users, err := database.Users(bson.M{
		"username": user.Username,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errRes{
			Status:  http.StatusInternalServerError,
			Message: "Internal server error",
		})
		return
	}
	if len(users) > 0 && *users[0].ID != oid {
		// category name exists
		c.JSON(http.StatusConflict, errRes{
			Status:  http.StatusConflict,
			Message: "Username already exists",
		})
		return
	}

	err = database.UpdateUser(
		bson.M{
			"_id": oid,
		},
		user,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errRes{
			Status:  http.StatusInternalServerError,
			Message: "Internal server error",
		})
		return
	}

	user.ID = &oid
	c.JSON(http.StatusCreated, user)
}

// UpdateUserPassword handles the PUT request
// for url path "/admin/user/:id/password"
func UpdateUserPassword(c *gin.Context) {
	type newPassword struct {
		Password string `json:"password" binding:"required"`
	}

	// parse object id from url path
	if !bson.IsObjectIdHex(c.Param("id")) {
		c.JSON(http.StatusBadRequest, errRes{
			Status:  http.StatusBadRequest,
			Message: "Invaild id",
		})
		return
	}

	oid := bson.ObjectIdHex(c.Param("id"))

	// get user id
	idStr, ok := c.Get("user_id")
	if !ok || !bson.IsObjectIdHex(idStr.(string)) {
		c.JSON(http.StatusUnauthorized, errRes{
			Status:  http.StatusUnauthorized,
			Message: "Invalid JWT token",
		})
		return
	}
	userID := bson.ObjectIdHex(idStr.(string))

	if userID != oid {
		c.JSON(http.StatusForbidden, errRes{
			Status:  http.StatusForbidden,
			Message: "Can't change other user's password",
		})
		return
	}

	newPsw := new(newPassword)
	if err := c.ShouldBindJSON(newPsw); err != nil {
		c.JSON(http.StatusBadRequest, errRes{
			Status:  http.StatusBadRequest,
			Message: "Bad Request",
		})
		return
	}

	var user structure.User
	var err error
	user.PasswordHash, err = bcrypt.GenerateFromPassword(
		[]byte(newPsw.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errRes{
			Status:  http.StatusInternalServerError,
			Message: "Internal server error",
		})
		return
	}

	err = database.UpdateUser(
		bson.M{
			"_id": oid,
		},
		user,
	)
	if err != nil {
		if err == database.ErrNoUser {
			c.JSON(http.StatusNotFound, errRes{
				Status:  http.StatusNotFound,
				Message: "No such user",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, errRes{
			Status:  http.StatusInternalServerError,
			Message: "Internal server error",
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
