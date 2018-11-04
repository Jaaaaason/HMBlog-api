package database

import (
	"context"
	"errors"
	"math/rand"

	"github.com/mongodb/mongo-go-driver/bson"

	"github.com/jaaaaason/hmblog/logger"
	"github.com/jaaaaason/hmblog/structure"

	"github.com/mongodb/mongo-go-driver/mongo"

	"golang.org/x/crypto/bcrypt"
)

// initBlogUser checks if there are blog users in database,
// create a default blog user "admin" with random password if not
func initBlogUser() error {
	ctx := context.Background()

	client, err := mongo.Connect(ctx, connectURI, nil)
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)

	c := client.Database(dbName).Collection("users")
	count, err := c.Count(ctx, nil, nil)
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

	_, err = c.InsertOne(ctx, bson.NewDocument(
		bson.EC.String("username", username),
		bson.EC.Binary("password_hash", pswHash),
	))
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
func User(ctx context.Context, filter map[string]interface{}) (*structure.User, error) {
	user := new(structure.User)

	client, err := mongo.Connect(ctx, connectURI, nil)
	if err != nil {
		return user, err
	}

	c := client.Database(dbName).Collection("users")

	d := bson.NewDocument()
	for key, val := range filter {
		d.Append(bson.EC.Interface(key, val))
	}

	err = c.FindOne(ctx, d, nil).Decode(user)
	if err != nil && err == mongo.ErrNoDocuments {
		return user, ErrNoUser
	}

	return user, err
}
