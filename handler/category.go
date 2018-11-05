package handler

import (
	"net/http"

	"gopkg.in/go-playground/validator.v8"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"github.com/jaaaaason/hmblog/database"
	"github.com/jaaaaason/hmblog/structure"
)

// GetCategories handle GET request for url path "/categories"
func GetCategories(c *gin.Context) {
	categories, err := database.Categories(nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errRes{
			Status:  http.StatusInternalServerError,
			Message: "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// GetCategory handle GET request for url path "/categories/:id"
func GetCategory(c *gin.Context) {
	// parse object id from url path
	if !bson.IsObjectIdHex(c.Param("id")) {
		c.JSON(http.StatusBadRequest, errRes{
			Status:  http.StatusBadRequest,
			Message: "Invaild id",
		})
		return
	}

	oid := bson.ObjectIdHex(c.Param("id"))

	categories, err := database.Categories(map[string]interface{}{
		"_id": oid,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errRes{
			Status:  http.StatusInternalServerError,
			Message: "Internal server error",
		})
		return
	}

	if len(categories) < 1 {
		c.JSON(http.StatusBadRequest, errRes{
			Status:  http.StatusBadRequest,
			Message: "No category found with id " + c.Param("id"),
		})
		return
	}

	c.JSON(http.StatusOK, categories[0])
}

// PostCategory handles POST request for url path "/admin/categories"
func PostCategory(c *gin.Context) {
	category := new(structure.Category)
	if err := c.ShouldBindJSON(category); err != nil {
		c.JSON(http.StatusBadRequest, errRes{
			Status:  http.StatusBadRequest,
			Message: "Bad request",
		})
		return
	}

	var err error
	category.ID, err = database.InsertCategory(category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errRes{
			Status:  http.StatusInternalServerError,
			Message: "Internal server error",
		})
		return
	}

	category.BlogCount = 0
	c.JSON(http.StatusCreated, category)
}

// UpdateCategory handles the PUT and PATCH request
// for url path "/admin/categories"
func UpdateCategory(c *gin.Context) {
	// parse object id from url path
	if !bson.IsObjectIdHex(c.Param("id")) {
		c.JSON(http.StatusBadRequest, errRes{
			Status:  http.StatusBadRequest,
			Message: "Invaild id",
		})
		return
	}

	oid := bson.ObjectIdHex(c.Param("id"))

	categories, err := database.Categories(map[string]interface{}{
		"_id": oid,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errRes{
			Status:  http.StatusInternalServerError,
			Message: "Internal server error",
		})
		return
	}

	if len(categories) < 1 {
		c.JSON(http.StatusBadRequest, errRes{
			Status:  http.StatusBadRequest,
			Message: "No category found with id " + c.Param("id"),
		})
		return
	}

	var category structure.Category
	if c.Request.Method == "PUT" {
		// for PUT request, use a new category struct,
		// binding with the request body, so the category
		// will be exactly the same as request body,
		// the value of some field that doesn't will be empty
		category = structure.Category{}
		err = c.ShouldBindJSON(&category)
		if err != nil {
			c.JSON(http.StatusBadRequest, errRes{
				Status:  http.StatusBadRequest,
				Message: "Bad request",
			})
			return
		}
	} else if c.Request.Method == "PATCH" {
		// for PATCH request, use the origin category just got before,
		// binding with request body, so the value of some field that
		// doesn't provide will not change
		category = categories[0]
		err = c.ShouldBindJSON(&category)
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

	err = database.UpdateCategory(oid, &category)
	if err != nil {
		if err == database.ErrNoCategory {
			c.JSON(http.StatusNotFound, errRes{
				Status:  http.StatusNotFound,
				Message: "No category found with id " + c.Param("id"),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, errRes{
			Status:  http.StatusInternalServerError,
			Message: "Internal server error",
		})
		return
	}

	category.ID = oid
	c.JSON(http.StatusCreated, category)
}
