package database

import (
	"errors"

	"github.com/globalsign/mgo"
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

// InsertPost inserts a post to database
func InsertPost(post *structure.Post) error {
	session := mgoSession.Copy()
	defer session.Close()

	c := session.DB(dbName).C("posts")

	if post.ID == nil {
		post.ID = new(bson.ObjectId)
	}
	*post.ID = bson.NewObjectId()

	return c.Insert(post)
}

// ErrNoPost returned when no category found
var ErrNoPost = errors.New("no such post")

// UpdatePost updates a exist post
func UpdatePost(filter bson.M, post structure.Post) error {
	session := mgoSession.Copy()
	defer session.Close()

	c := session.DB(dbName).C("posts")

	err := c.Update(
		filter,
		bson.M{
			"$set": post,
		},
	)
	if err != nil && err == mgo.ErrNotFound {
		return ErrNoPost
	}

	return err
}
