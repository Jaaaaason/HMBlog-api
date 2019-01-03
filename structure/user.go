package structure

import "github.com/globalsign/mgo/bson"

// User the blog user struct
type User struct {
	ID           *bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Username     string         `json:"username" bson:"username,omitempty"`
	PasswordHash []byte         `json:"-" bson:"password_hash,omitempty"`
}
