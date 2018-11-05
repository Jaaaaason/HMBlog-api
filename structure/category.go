package structure

import "github.com/globalsign/mgo/bson"

// Category the blog category struct
type Category struct {
	ID        bson.ObjectId `json:"id" bson:"-"`
	Name      string        `json:"name" bson:"name" binding:"required"`
	BlogCount int           `json:"blog_count" bson:"-"`
}
