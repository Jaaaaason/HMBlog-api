package handler

import (
	"net/http"
	"strings"

	"gopkg.in/go-playground/validator.v8"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"github.com/jaaaaason/hmblog/database"
	"github.com/jaaaaason/hmblog/structure"
)

// GetCategories handles GET request for url path "/categories"
func GetCategories(c *gin.Context) {
	categories, err := database.Categories(nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errRes{
			Status:  http.StatusInternalServerError,
			Message: "Internal server error",
		})
		return
	}

	for i := range categories {
		categories[i].PostCount, err = database.PostCount(
			bson.M{
				"category_id": categories[i].ID,
				"is_publish":  true,
			},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errRes{
				Status:  http.StatusInternalServerError,
				Message: "Internal server error",
			})
			return
		}
	}

	c.JSON(http.StatusOK, categories)
}

// GetCategory handles GET request for url path "/categories/:id"
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

	categories, err := database.Categories(bson.M{
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

	categories[0].PostCount, err = database.PostCount(
		bson.M{
			"category_id": categories[0].ID,
			"is_publish":  true,
		},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errRes{
			Status:  http.StatusInternalServerError,
			Message: "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, categories[0])
}

// GetAdminCategories handles GET request for url path "/admin/categories"
func GetAdminCategories(c *gin.Context) {
	categories, err := database.Categories(nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errRes{
			Status:  http.StatusInternalServerError,
			Message: "Internal server error",
		})
		return
	}

	idStr, ok := c.Get("user_id")
	if !ok || !bson.IsObjectIdHex(idStr.(string)) {
		c.JSON(http.StatusUnauthorized, errRes{
			Status:  http.StatusUnauthorized,
			Message: "Invalid JWT token",
		})
		return
	}
	userID := bson.ObjectIdHex(idStr.(string))

	for i := range categories {
		// count published post or
		// unpublish post that belongs to current user
		categories[i].PostCount, err = database.PostCount(
			bson.M{
				"$or": []bson.M{
					bson.M{
						"category_id": categories[i].ID,
						"is_publish":  true,
					},
					bson.M{
						"category_id": categories[i].ID,
						"is_publish":  false,
						"user_id":     userID,
					},
				},
			},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errRes{
				Status:  http.StatusInternalServerError,
				Message: "Internal server error",
			})
			return
		}
	}

	c.JSON(http.StatusOK, categories)
}

// GetAdminCategory handles GET request for url path "/admin/categories/:id"
func GetAdminCategory(c *gin.Context) {
	// parse object id from url path
	if !bson.IsObjectIdHex(c.Param("id")) {
		c.JSON(http.StatusBadRequest, errRes{
			Status:  http.StatusBadRequest,
			Message: "Invaild id",
		})
		return
	}

	oid := bson.ObjectIdHex(c.Param("id"))

	idStr, ok := c.Get("user_id")
	if !ok || !bson.IsObjectIdHex(idStr.(string)) {
		c.JSON(http.StatusUnauthorized, errRes{
			Status:  http.StatusUnauthorized,
			Message: "Invalid JWT token",
		})
		return
	}
	userID := bson.ObjectIdHex(idStr.(string))

	categories, err := database.Categories(bson.M{
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

	// count published post or
	// unpublish post that belongs to current user
	categories[0].PostCount, err = database.PostCount(
		bson.M{
			"$or": []bson.M{
				bson.M{
					"category_id": categories[0].ID,
					"is_publish":  true,
				},
				bson.M{
					"category_id": categories[0].ID,
					"is_publish":  false,
					"user_id":     userID,
				},
			},
		},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errRes{
			Status:  http.StatusInternalServerError,
			Message: "Internal server error",
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

	// trim space
	category.Name = strings.TrimSpace(category.Name)
	if category.Name == "" {
		// empty category name
		c.JSON(http.StatusBadRequest, errRes{
			Status:  http.StatusBadRequest,
			Message: "Category name shouldn't be just some whitespace",
		})
		return
	}

	categories, err := database.Categories(bson.M{
		"name": category.Name,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errRes{
			Status:  http.StatusInternalServerError,
			Message: "Internal server error",
		})
		return
	}

	if len(categories) > 0 {
		// category name exists
		c.JSON(http.StatusConflict, errRes{
			Status:  http.StatusConflict,
			Message: "Category name already exists",
		})
		return
	}

	err = database.InsertCategory(category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errRes{
			Status:  http.StatusInternalServerError,
			Message: "Internal server error",
		})
		return
	}

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

	categories, err := database.Categories(bson.M{
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

	// set field ID empty value to omit it
	category.ID = nil

	// trim space
	category.Name = strings.TrimSpace(category.Name)
	if category.Name == "" {
		// empty category name
		c.JSON(http.StatusBadRequest, errRes{
			Status:  http.StatusBadRequest,
			Message: "Category name shouldn't be just some whitespace",
		})
		return
	}

	categories, err = database.Categories(bson.M{
		"name": category.Name,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errRes{
			Status:  http.StatusInternalServerError,
			Message: "Internal server error",
		})
		return
	}
	if len(categories) > 0 && *categories[0].ID != oid {
		// category name exists
		c.JSON(http.StatusConflict, errRes{
			Status:  http.StatusConflict,
			Message: "Category name already exists",
		})
		return
	}

	err = database.UpdateCategories(
		bson.M{
			"_id": oid,
		},
		category,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errRes{
			Status:  http.StatusInternalServerError,
			Message: "Internal server error",
		})
		return
	}

	category.ID = &oid
	c.JSON(http.StatusCreated, category)
}

// DeleteCategory handles the DELETE request
// of url path "/admin/categories/:id"
func DeleteCategory(c *gin.Context) {
	// parse object id from url path
	if !bson.IsObjectIdHex(c.Param("id")) {
		c.JSON(http.StatusBadRequest, errRes{
			Status:  http.StatusBadRequest,
			Message: "Invaild id",
		})
		return
	}
	oid := bson.ObjectIdHex(c.Param("id"))

	err := database.RemoveCategories(bson.M{
		"_id": oid,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errRes{
			Status:  http.StatusInternalServerError,
			Message: "Internal server error",
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
