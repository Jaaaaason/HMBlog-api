package configer

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// Configer the configuration struct
type Configer struct {
	MongoDBListen int    `json:"mongodb_listen"`
	DBName        string `json:"database_name"`
	LogFile       string `json:"log_file"`
}

// Config the global config
var Config Configer

// Initialize reads the config file and
// stores the config into the variable Config
func Initialize(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &Config)

	return err
}
