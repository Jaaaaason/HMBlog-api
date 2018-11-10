package database

import (
	"errors"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"

	"github.com/jaaaaason/hmblog/structure"
)

// ErrNoCategory returned when no category found
var ErrNoCategory = errors.New("no such category")

// Categories returns all categories which match the filter
func Categories(filter bson.M) ([]structure.Category, error) {
	var categories []structure.Category

	session := mgoSession.Copy()
	defer session.Close()

	c := session.DB(dbName).C("categories")

	pipeline := []bson.M{
		bson.M{
			"$match": filter,
		},
		bson.M{
			"$project": bson.M{
				"_id":  1,
				"name": 1,
			},
		},
	}

	err := c.Pipe(pipeline).All(&categories)

	return categories, err
}

// InsertCategory inserts a category
func InsertCategory(category *structure.Category) error {
	session := mgoSession.Copy()
	defer session.Close()

	c := session.DB(dbName).C("categories")

	if category.ID == nil {
		category.ID = new(bson.ObjectId)
	}
	*category.ID = bson.NewObjectId()

	return c.Insert(category)
}

// UpdateCategories updates all categories that matches the filter,
// ErrNoCategory returned when destination category doesn't exist
func UpdateCategories(filter bson.M, category structure.Category) error {
	session := mgoSession.Copy()
	defer session.Close()

	c := session.DB(dbName).C("categories")

	_, err := c.UpdateAll(
		filter,
		bson.M{
			"$set": category,
		},
	)
	if err != nil && err == mgo.ErrNotFound {
		return ErrNoCategory
	}

	return err
}

// RemoveCategories removes all categories that matches the filter
func RemoveCategories(filter bson.M) error {
	session := mgoSession.Copy()
	defer session.Close()

	c := session.DB(dbName).C("categories")

	_, err := c.RemoveAll(filter)
	return err
}
