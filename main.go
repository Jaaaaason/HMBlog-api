package main

import (
	"flag"

	"github.com/jaaaaason/hmblog/configer"
	"github.com/jaaaaason/hmblog/database"
	"github.com/jaaaaason/hmblog/logger"
)

func main() {
	// get config file's path with commandline arg
	confFilepath := flag.String("c", "", "the config file's path")

	var err error
	if *confFilepath == "" {
		// use default config file's path
		err = configer.Initialize("config.json")
	} else {
		err = configer.Initialize(*confFilepath)
	}

	if err != nil {
		logger.Fatal(err.Error())
	}

	// desire log file is given
	if configer.Config.LogFile != "" {
		err = logger.SetOutputFile(configer.Config.LogFile)
		if err != nil {
			logger.Fatal(err.Error())
		}
	}

	// initialize the database's connection
	err = database.Initialize()
	if err != nil {
		logger.Fatal(err.Error())
	}
}
