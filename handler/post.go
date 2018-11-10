package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"github.com/jaaaaason/hmblog/database"
	"github.com/jaaaaason/hmblog/structure"
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

// PostPost handles the POST request of url path "/admin/posts"
func PostPost(c *gin.Context) {
	post := new(structure.Post)
	if err := c.ShouldBindJSON(post); err != nil {
		c.JSON(http.StatusBadRequest, errRes{
			Status:  http.StatusBadRequest,
			Message: "Bad request",
		})
		return
	}

	// trim space
	post.Title = strings.TrimSpace(post.Title)
	if post.Title == "" {
		// empty category name
		c.JSON(http.StatusBadRequest, errRes{
			Status:  http.StatusBadRequest,
			Message: "Title shouldn't be just some whitespace",
		})
		return
	}

	posts, err := database.Posts(bson.M{
		"title": post.Title,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errRes{
			Status:  http.StatusInternalServerError,
			Message: "Internal server error",
		})
		return
	}

	if len(posts) > 0 {
		// post title exists
		c.JSON(http.StatusConflict, errRes{
			Status:  http.StatusConflict,
			Message: "Post with this title already exists",
		})
		return
	}

	// get user id
	idStr, ok := c.Get("user_id")
	if !ok || !bson.IsObjectIdHex(idStr.(string)) {
		c.JSON(http.StatusUnauthorized, errRes{
			Status:  http.StatusUnauthorized,
			Message: "Invalid JWT token",
		})
		return
	}
	post.UserID = new(bson.ObjectId)
	*post.UserID = bson.ObjectIdHex(idStr.(string))

	post.CategoryName = strings.TrimSpace(post.CategoryName)
	if post.CategoryName != "" {
		categories, err := database.Categories(bson.M{
			"name": post.CategoryName,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, errRes{
				Status:  http.StatusInternalServerError,
				Message: "Internal server error",
			})
			return
		}

		if len(categories) > 0 {
			post.CategoryID = categories[0].ID
			categories[0].PostCount, _ = database.PostCount(bson.M{
				"$or": []bson.M{
					bson.M{
						"category_id": categories[0].ID,
						"is_publish":  true,
					},
					bson.M{
						"category_id": categories[0].ID,
						"is_publish":  false,
						"user_id":     post.UserID,
					},
				},
			})
			categories[0].PostCount++
			post.Category = &categories[0]
		} else {
			category := structure.Category{
				Name: post.CategoryName,
			}
			err = database.InsertCategory(&category)
			if err != nil {
				c.JSON(http.StatusInternalServerError, errRes{
					Status:  http.StatusInternalServerError,
					Message: "Internal server error",
				})
				return
			}

			post.CategoryID = category.ID
			post.Category = &category
			post.Category.PostCount = 1
		}
	}

	post.CreatedAt = time.Now()
	post.UpdatedAt = post.CreatedAt

	err = database.InsertPost(post)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errRes{
			Status:  http.StatusInternalServerError,
			Message: "Internal server error",
		})
		return

		// TODO: delete category if it is just created
	}

	// retrieve user
	post.User = new(structure.User)
	*post.User, _ = database.User(bson.M{
		"_id": post.UserID,
	})

	c.JSON(http.StatusCreated, post)
}
