package database

import (
	"fmt"
	"strings"

	"github.com/globalsign/mgo"

	"github.com/jaaaaason/hmblog/configer"
)

var mgoSession *mgo.Session // original mgo session
var dbName = "blog"

// Initialize initializes the connection uri
// and test the connection
func Initialize() error {
	connectURI := fmt.Sprintf("mongodb://%s:%d",
		configer.Config.MongoDBHost,
		configer.Config.MongoDBListen,
	)

	if strings.TrimSpace(configer.Config.DBName) != "" {
		// database name is given
		dbName = strings.TrimSpace(configer.Config.DBName)
	}

	// test connection
	var err error
	mgoSession, err = mgo.Dial(connectURI)
	if err != nil {
		return err
	}

	return initBlogUser()
}

// CloseSession closes the original mgo session "mgoSession"
func CloseSession() {
	mgoSession.Close()
}
