package database

import (
	"errors"
	"math/rand"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"

	"github.com/jaaaaason/hmblog/logger"
	"github.com/jaaaaason/hmblog/structure"

	"golang.org/x/crypto/bcrypt"
)

// initBlogUser checks if there are blog users in database,
// create a default blog user "admin" with random password if not
func initBlogUser() error {
	session := mgoSession.Copy()
	defer session.Close()

	c := session.DB(dbName).C("users")
	count, err := c.Count()
	if err != nil {
		return err
	}

	// has user, no need to create a default one
	if count > 0 {
		return nil
	}

	username := "admin"
	psw := randomPassword() // create random password
	pswHash, err := bcrypt.GenerateFromPassword(psw, bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = c.Insert(bson.M{
		"username":      username,
		"password_hash": pswHash,
	})
	if err != nil {
		return err
	}

	logger.Info("created a default user:\n" +
		"username: " + username + "\n" +
		"password: " + string(psw))

	return nil
}

// randomPassword creates a 12-length random password
func randomPassword() []byte {
	characters := []byte("0123456789" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz")

	psw := make([]byte, 12)
	for i := 0; i < 12; i++ {
		psw[i] = characters[rand.Intn(62)]
	}

	return psw
}

// ErrNoUser returned when no user found
var ErrNoUser = errors.New("no such user")

// User returns one user that match the filter
func User(filter map[string]interface{}) (structure.User, error) {
	var user structure.User

	session := mgoSession.Copy()
	defer session.Close()

	c := session.DB(dbName).C("users")

	err := c.Find(filter).One(&user)
	if err != nil && err == mgo.ErrNotFound {
		return user, ErrNoUser
	}

	return user, err
}
