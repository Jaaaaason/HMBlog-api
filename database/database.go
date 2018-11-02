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

// Initialize initializes the connection uri
// and test the connection
func Initialize() error {
	connectURI = fmt.Sprintf("mongodb://%s:%s@%s:%d/",
		configer.Config.DBUser,
		configer.Config.DBUserPassword,
		configer.Config.MongoDBHost,
		configer.Config.MongoDBListen,
	)

	if strings.TrimSpace(configer.Config.DBName) != "" {
		// database name is given
		connectURI += strings.TrimSpace(configer.Config.DBName)
	} else {
		// use default database name "blog"
		connectURI += "blog"
	}

	client, err := mongo.NewClient(connectURI)
	if err != nil {
		return err
	}

	// test connection
	return client.Connect(context.Background())
}
