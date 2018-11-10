package structure

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

// Post the blog post struct
type Post struct {
	ID           *bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Title        string         `json:"title" bson:"title,omitempty" binding:"required"`
	Content      string         `json:"content" bson:"content,omitempty" binding:"required"`
	IsPublish    *bool          `json:"is_publish" bson:"is_publish,omitempty" binding:"required"`
	CategoryID   *bson.ObjectId `json:"-" bson:"category_id"`
	Category     *Category      `json:"category" bson:"-"`
	CategoryName string         `json:"category_name" bson:"-"`
	Tags         []string       `json:"tags" bson:"tags"`
	UserID       *bson.ObjectId `json:"-" bson:"user_id"`
	User         *User          `json:"user" bson:"-"`
	CreatedAt    time.Time      `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at" bson:"updated_at"`
}
