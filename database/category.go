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
func Categories(filter map[string]interface{}) ([]structure.Category, error) {
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
				"blog_count": bson.M{
					"$size": "$blog_ids",
				},
			},
		},
	}

	err := c.Pipe(pipeline).All(&categories)

	return categories, err
}

// InsertCategory inserts a category
func InsertCategory(category *structure.Category) (*bson.ObjectId, error) {
	id := bson.NewObjectId()

	session := mgoSession.Copy()
	defer session.Close()

	c := session.DB(dbName).C("categories")

	err := c.Insert(bson.M{
		"_id":      id,
		"name":     category.Name,
		"blog_ids": []bson.ObjectId{},
	})
	if err != nil {
		return &id, err
	}

	return &id, nil
}

// UpdateCategory updates a exists category with given category,
// ErrNoCategory returned when destination category doesn't exist
func UpdateCategory(id bson.ObjectId, category *structure.Category) error {
	session := mgoSession.Copy()
	defer session.Close()

	c := session.DB(dbName).C("categories")

	err := c.Update(
		bson.M{
			"_id": id,
		},
		bson.M{
			"$set": category,
		},
	)
	if err != nil {
		if err == mgo.ErrNotFound {
			return ErrNoCategory
		}

		return err
	}

	return nil
}
