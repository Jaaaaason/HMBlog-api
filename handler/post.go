package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"github.com/jaaaaason/hmblog/database"
)

// GetPosts handles the GET request of url path "/posts"
func GetPosts(c *gin.Context) {
	posts, err := database.Posts(bson.M{
		"is_publish": true,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errRes{
			Status:  http.StatusInternalServerError,
			Message: "Internal server error",
		})
		return
	}

	for i := range posts {
		if posts[i].CategoryID != nil {
			// retrieve post's category
			categories, err := database.Categories(bson.M{
				"_id": posts[i].CategoryID,
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, errRes{
					Status:  http.StatusInternalServerError,
					Message: "Internal server error",
				})
				return
			}
			if len(categories) > 0 {
				categories[0].PostCount, err = database.PostCount(bson.M{
					"category_id": categories[0].ID,
					"is_publish":  true,
				})
				if err != nil {
					c.JSON(http.StatusInternalServerError, errRes{
						Status:  http.StatusInternalServerError,
						Message: "Internal server error",
					})
					return
				}

				posts[i].Category = &categories[0]
			}
		}

		if posts[i].UserID != nil {
			// retrieve post's owner
			user, err := database.User(bson.M{
				"_id": posts[i].UserID,
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, errRes{
					Status:  http.StatusInternalServerError,
					Message: "Internal server error",
				})
				return
			}
			posts[i].User = &user
		}
	}

	c.JSON(http.StatusOK, posts)
}

// GetPost handles the GET request of url path "/posts/:id"
func GetPost(c *gin.Context) {
	// parse object id from url path
	if !bson.IsObjectIdHex(c.Param("id")) {
		c.JSON(http.StatusBadRequest, errRes{
			Status:  http.StatusBadRequest,
			Message: "Invaild id",
		})
		return
	}

	oid := bson.ObjectIdHex(c.Param("id"))

	posts, err := database.Posts(bson.M{
		"_id":        oid,
		"is_publish": true,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errRes{
			Status:  http.StatusInternalServerError,
			Message: "Internal server error",
		})
		return
	}

	if len(posts) < 1 {
		c.JSON(http.StatusBadRequest, errRes{
			Status:  http.StatusBadRequest,
			Message: "No post found",
		})
		return
	}

	if posts[0].CategoryID != nil {
		// retrieve post's category
		categories, err := database.Categories(bson.M{
			"_id": posts[0].CategoryID,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, errRes{
				Status:  http.StatusInternalServerError,
				Message: "Internal server error",
			})
			return
		}
		if len(categories) > 0 {
			categories[0].PostCount, err = database.PostCount(bson.M{
				"category_id": categories[0].ID,
				"is_publish":  true,
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, errRes{
					Status:  http.StatusInternalServerError,
					Message: "Internal server error",
				})
				return
			}

			posts[0].Category = &categories[0]
		}
	}

	if posts[0].UserID != nil {
		// retrieve post's owner
		user, err := database.User(bson.M{
			"_id": posts[0].UserID,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, errRes{
				Status:  http.StatusInternalServerError,
				Message: "Internal server error",
			})
			return
		}
		posts[0].User = &user
	}

	c.JSON(http.StatusOK, posts[0])
}

// GetAdminPosts handles the GET request of url path "/admin/posts"
func GetAdminPosts(c *gin.Context) {
	idStr, ok := c.Get("user_id")
	if !ok || !bson.IsObjectIdHex(idStr.(string)) {
		c.JSON(http.StatusUnauthorized, errRes{
			Status:  http.StatusUnauthorized,
			Message: "Invalid JWT token",
		})
		return
	}
	userID := bson.ObjectIdHex(idStr.(string))

	posts, err := database.Posts(bson.M{
		"$or": []bson.M{
			bson.M{
				"is_publish": true,
			},
			bson.M{
				"is_publish": false,
				"user_id":    userID,
			},
		},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errRes{
			Status:  http.StatusInternalServerError,
			Message: "Internal server error",
		})
		return
	}

	for i := range posts {
		if posts[i].CategoryID != nil {
			// retrieve post's category
			categories, err := database.Categories(bson.M{
				"_id": posts[i].CategoryID,
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, errRes{
					Status:  http.StatusInternalServerError,
					Message: "Internal server error",
				})
				return
			}
			if len(categories) > 0 {
				categories[0].PostCount, err = database.PostCount(bson.M{
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
				})
				if err != nil {
					c.JSON(http.StatusInternalServerError, errRes{
						Status:  http.StatusInternalServerError,
						Message: "Internal server error",
					})
					return
				}

				posts[i].Category = &categories[0]
			}
		}

		if posts[i].UserID != nil {
			// retrieve post's owner
			user, err := database.User(bson.M{
				"_id": posts[i].UserID,
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, errRes{
					Status:  http.StatusInternalServerError,
					Message: "Internal server error",
				})
				return
			}
			posts[i].User = &user
		}
	}

	c.JSON(http.StatusOK, posts)
}

//
func GetAdminPost(c *gin.Context) {
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

	posts, err := database.Posts(bson.M{
		"$or": []bson.M{
			bson.M{
				"_id":        oid,
				"is_publish": true,
			},
			bson.M{
				"_id":        oid,
				"is_publish": false,
				"user_id":    userID,
			},
		},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errRes{
			Status:  http.StatusInternalServerError,
			Message: "Internal server error",
		})
		return
	}

	if len(posts) < 1 {
		c.JSON(http.StatusBadRequest, errRes{
			Status:  http.StatusBadRequest,
			Message: "No post found",
		})
		return
	}

	if posts[0].CategoryID != nil {
		// retrieve post's category
		categories, err := database.Categories(bson.M{
			"_id": posts[0].CategoryID,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, errRes{
				Status:  http.StatusInternalServerError,
				Message: "Internal server error",
			})
			return
		}
		if len(categories) > 0 {
			categories[0].PostCount, err = database.PostCount(bson.M{
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
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, errRes{
					Status:  http.StatusInternalServerError,
					Message: "Internal server error",
				})
				return
			}

			posts[0].Category = &categories[0]
		}
	}

	if posts[0].UserID != nil {
		// retrieve post's owner
		user, err := database.User(bson.M{
			"_id": posts[0].UserID,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, errRes{
				Status:  http.StatusInternalServerError,
				Message: "Internal server error",
			})
			return
		}
		posts[0].User = &user
	}

	c.JSON(http.StatusOK, posts[0])
}
