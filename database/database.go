package database

import (
	"context"
	"fmt"
	"strings"

	"github.com/jaaaaason/hmblog/configer"

	"github.com/mongodb/mongo-go-driver/mongo"
)

// connectURI the connection uri for mongodb
var connectURI string
var dbName = "blog"

// Initialize initializes the connection uri
// and test the connection
func Initialize() error {
	connectURI = fmt.Sprintf("mongodb://%s:%s@%s:%d",
		configer.Config.DBUser,
		configer.Config.DBUserPassword,
		configer.Config.MongoDBHost,
		configer.Config.MongoDBListen,
	)

	if strings.TrimSpace(configer.Config.DBName) != "" {
		// database name is given
		dbName = strings.TrimSpace(configer.Config.DBName)
	}

	// test connection
	_, err := mongo.Connect(context.Background(), connectURI, nil)
	if err != nil {
		return err
	}

	return initBlogUser()
}
