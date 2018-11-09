package database

import "github.com/globalsign/mgo/bson"

// PostCount returns the amount of post that matches the filter
func PostCount(filter bson.M) (int, error) {
	session := mgoSession.Copy()
	defer session.Close()

	c := session.DB(dbName).C("posts")

	return c.Find(filter).Count()
}
