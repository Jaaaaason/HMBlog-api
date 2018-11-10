package database

import (
	"github.com/globalsign/mgo/bson"
	"github.com/jaaaaason/hmblog/structure"
)

// PostCount returns the amount of post that matches the filter
func PostCount(filter bson.M) (int, error) {
	session := mgoSession.Copy()
	defer session.Close()

	c := session.DB(dbName).C("posts")

	return c.Find(filter).Count()
}

// Posts retrieves posts that matches the filter from database
func Posts(filter bson.M) ([]structure.Post, error) {
	session := mgoSession.Copy()
	defer session.Close()

	c := session.DB(dbName).C("posts")

	var posts []structure.Post
	err := c.Find(filter).All(&posts)

	return posts, err
}
