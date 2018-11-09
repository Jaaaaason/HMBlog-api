package structure

import "github.com/globalsign/mgo/bson"

// Category the blog category struct
type Category struct {
	ID        *bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name      string         `json:"name" bson:"name" binding:"required"`
	PostCount int            `json:"post_count" bson:"-"`
}
